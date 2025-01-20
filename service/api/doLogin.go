package api

import (
	"crypto/rand"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"math/big"
	"net/http"
)

func generateAPIKey() (string, error) {
	const apiKeyLength = 16
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var apiKey []byte

	for i := 0; i < apiKeyLength; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "error key", err
		}
		apiKey = append(apiKey, charset[index.Int64()])
	}

	return string(apiKey), nil
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var requestBody struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := requestBody.Name
	if len(username) < 3 || len(username) > 16 {
		json.NewEncoder(w).Encode(map[string]string{"error": "Username must be between 3 and 16 characters"})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := rt.db.CheckIfUserExists(username)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var apiKey string
	var userID int
	if !exists {
		apiKey, err = generateAPIKey()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Error generating API key"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userID, err = rt.db.AddNewUser(username, apiKey)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Error adding new user"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		userID, err = rt.db.GetUserID(username)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving user ID"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		apiKey, err = rt.db.GetUserKey(userID)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving API key"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	response := struct {
		Username string `json:"username"`
		UserID   int    `json:"userId"`
		APIKey   string `json:"apiKey"`
	}{
		Username: username,
		UserID:   userID,
		APIKey:   apiKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
