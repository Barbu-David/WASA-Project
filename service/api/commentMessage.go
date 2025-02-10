package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	conv_id_param := ps.ByName("convid")

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
		return
		}

	// And be a member of the conversation

	member, err := rt.db.IsMemberConversation(user_id, conv_id)

	if err != nil || member != true {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Auth error"})
		return
	}

        m_id_param := ps.ByName("MessageId")

        m_id, err := strconv.Atoi(m_id_param)
        if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                _ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid mid"})
                return
        }

        var requestBody struct {
                comment string `json:"comment"`
        }

        if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
                w.WriteHeader(http.StatusBadRequest)
                ctx.Logger.WithError(err).Error("Encoding failed fail")
                _ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
                return
        }

        if len(requestBody.comment) > 10000 {
                w.WriteHeader(http.StatusBadRequest)
                _ = json.NewEncoder(w).Encode(map[string]string{"error": "Comment too long"})
                return
        }

	//DB call

	err = rt.db.AddComment(user_id, m_id, requestBody.comment)
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "InternalServerError"})
		 ctx.Logger.WithError(err).Error("database fail")
		return 
	}

	w.WriteHeader(http.StatusNoContent)
}
