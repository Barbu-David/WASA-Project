package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) startConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Validate the Authorization header for a Bearer token.
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
		return
	}

	// Decode the JSON request body into a struct.
	var reqBody struct {
		UserIds []int `json:"userIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON provided"})
		return
	}

	// Check that the required field is present.
	if reqBody.UserIds == nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing required field: userIds"})
		return
	}

	// Check that the list does not exceed the maximum allowed size.
	if len(reqBody.UserIds) > 10000 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Too many user IDs provided"})
		return
	}

	var conv_id int
	var err error

	if len(reqBody.UserIds) > 1 {

		name0, err := rt.db.GetUserName(reqBody.UserIds[0])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}

		name1, err := rt.db.GetUserName(reqBody.UserIds[1])

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}

		conv_id, err = rt.db.NewConversation("Conversation between "+name0+" and "+name1, false)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}

	} else {
		conv_id, err = rt.db.NewConversation("Group conversation", true)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create conversation"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}
	}

	for _, userId := range reqBody.UserIds {

		err := rt.db.NewConversationMember(userId, conv_id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
			ctx.Logger.WithError(err).Error("Database fail")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conv_id)
}
