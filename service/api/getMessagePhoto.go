package api

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"image/gif"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) getMessagePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid message id"})
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

	m_id_param := ps.ByName("messageId")

	m_id, err := strconv.Atoi(m_id_param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid mid"})
		return
	}

	// GET THE GIF PHOTO

	photo, err := rt.db.GetMessagePhoto(m_id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Message not found"})
		ctx.Logger.WithError(err).Error("Database fail")
		return
	}

	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, photo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding photo"})
		ctx.Logger.WithError(err).Error("Encoding failed")
		return
	}

	w.Header().Set("Content-Type", "image/gif")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		ctx.Logger.WithError(err).Error("Error writing response")
	}

}
