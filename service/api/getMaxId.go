package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) getMaxId(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	authHeader := r.Header.Get("Authorization")
	fmt.Println("Header: ", authHeader, "\n")

	const bearerPrefix = " Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization format"})
		return
	}

	// Extract the token by removing 'Bearer ' prefix
	token := authHeader[len(bearerPrefix):]

	if token == "" { //The user should be able to get the name of any user, he just needs to have a valid token
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Empty token"})
	}

	maxID, err := rt.db.GetMaxUserID()
	if err != nil {
		// Log the error and return an internal server error response
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to retrieve maximum user ID"})
		ctx.Logger.WithError(err).Error("Database error while retrieving max user ID")
		return
	}

	// Respond with the maximum user ID
	response := struct {
		MaxUserID int `json:"maxUserId"`
	}{
		MaxUserID: maxID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
