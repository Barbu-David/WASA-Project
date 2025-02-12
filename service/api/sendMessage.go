package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"image/gif"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
	"wasatext/service/globaltime"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	user_id, err := rt.db.GetUserIDbyKey(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Auth error"})
		return
	}

	member, err := rt.db.IsMemberConversation(user_id, conv_id)
	if err != nil || !member {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Auth error"})
		return
	}

	contentType := r.Header.Get("Content-Type")
	var (
		messageText string
		isPhoto     bool
		photoData   *gif.GIF
	)

	switch contentType {
	case "application/json":
		var requestBody struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("JSON decode failed")
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
			return
		}
		messageText = requestBody.Message
		isPhoto = false
		photoData = nil

		if len(messageText) > 10000 {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Message too long"})
			return
		}

	case "image/gif":
		img, err := gif.DecodeAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("GIF decode failed")
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid GIF data"})
			return
		}
		messageText = ""
		isPhoto = true
		photoData = img

	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unsupported content type"})
		return
	}

	err = rt.db.SendMessage(user_id, conv_id, messageText, false, globaltime.Now(), isPhoto, photoData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "InternalServerError"})
		ctx.Logger.WithError(err).Error("database fail")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
