package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	userIDParam := ps.ByName("id")

	// Convert userID to integer
	requestedUserID, err := strconv.Atoi(userIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or unauthorized user ID"})
		return
	}

	actualKey, err := rt.db.GetUserKey(requestedUserID)

	// Check if the user and the security key match
	if actualKey != token {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	// Parse the new username from the request body
	var requestBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if len(requestBody.Name) < 3 || len(requestBody.Name) > 16 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Username must be between 3 and 16 characters"})
		return
	}
	
	err = rt.db.SetUserName(requestedUserID, requestBody.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update username"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
