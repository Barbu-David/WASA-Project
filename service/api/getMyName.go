package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the security token from the request header
	securityToken := r.Header.Get("security_token")
	if securityToken == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing security token"})
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Extract the user ID from the URL path
	userIDParam := ps.ByName("id")

	// Convert userID to integer
	requestedUserID, err := strconv.Atoi(userIDParam)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or unauthorized user ID"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve the user's name from the database
	username, err := rt.db.GetUserName(requestedUserID)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Respond with the username
	response := struct {
		Username string `json:"username"`
	}{
		Username: username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
