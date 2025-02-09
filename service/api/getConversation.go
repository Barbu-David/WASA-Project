package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization format"})
		return
	}

	token := authHeader[len(bearerPrefix):]

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Empty token"})
	}

	// Extract the user ID from the URL path
	conv_id_param := ps.ByName("convid")

	// Convert userID to integer
	conv_id, err := strconv.Atoi(conv_id_param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid conversation id"})
		return
	}

	// User must exist

	user_id, err := rt.db.GetUserIDbyKey(token)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Auth error"})
	}

	// And be a member of the conversation

	member, err := rt.db.IsMemberConversation(user_id, conv_id)

	if err != nil || member != true {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Auth error"})
	}

	participants, err := rt.db.GetConversationUsers(conv_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error retrieving conversation users")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "InternalServerError"})
		return
	}

	response := struct {
		Participants []int `json:"participants"`
		Messages     []int `json:"messages"`
	}{
		Participants: participants,
		Messages:     []int{},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
		ctx.Logger.WithError(err).Error("Encoding failed")
	}
}
