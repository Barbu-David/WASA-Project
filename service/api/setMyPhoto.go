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

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	actualKey, err := rt.db.GetUserKey(requestedUserID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database fail")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal error"})
		return
	}

	// Check if the user and the security key match
	if actualKey != token {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
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

	// Read the request body
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

	// Decode the GIF image
	gifImg, err := gif.DecodeAll(bytes.NewReader(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid GIF image"})
		return
	}

	// Update the user's photo in the database
	err = rt.db.SetUserPhoto(requestedUserID, gifImg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error setting user photo")
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
