package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
	"image/gif"
	"bytes"
)

func (rt *_router) getMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	// Extract the user ID from the URL path
	userIDParam := ps.ByName("id")

	// Convert userID to integer
	requestedUserID, err := strconv.Atoi(userIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or unauthorized user ID"})
		return
	}


	photo, err := rt.db.GetUserPhoto(requestedUserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
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
