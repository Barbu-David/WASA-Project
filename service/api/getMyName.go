package api

import (
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
	"strconv"
)

func (rt *_router) getMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract the user ID from the URL path
	userIDParam := ps.ByName("id")

	// Convert userID to integer
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		// Invalid user ID format
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user ID"})
		return
	}

	// Retrieve the user's name from the database
	username, err := rt.db.GetUserName(userID)
	if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
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
	json.NewEncoder(w).Encode(response)
}
