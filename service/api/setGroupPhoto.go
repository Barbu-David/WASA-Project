package api

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"image/gif"
	"io"
	"mime"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Check Content-Type is image/gif
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Missing Content-Type header"})
		return
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid Content-Type header"})
		return
	}

	if mediaType != "image/gif" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unsupported media type, expected image/gif"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Validate image size constraints
	if len(body) < 40 || len(body) > 10000000 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Image size must be between 40 and 10,000,000 bytes"})
		return
	}

	group, err := rt.db.IsGroupConversation(conv_id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to look up is_group")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	if group != true {
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Not a group"})
		return
	}

	// Decode the GIF image
	gifImg, err := gif.DecodeAll(bytes.NewReader(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid GIF image"})
		return
	}

	err = rt.db.SetConversationPhoto(conv_id, gifImg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error setting user photo")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
