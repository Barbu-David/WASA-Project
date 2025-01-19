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
			return "", err
		}
		apiKey = append(apiKey, charset[index.Int64()])
	}

	return string(apiKey), nil
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var requestBody struct {
		Name string `json:"name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&requestBody)

	username := requestBody.Name

	exists, _ := rt.db.CheckIfUserExists(username)

	var apiKey string
	if !exists {
		apiKey, _ = generateAPIKey()

		rt.db.AddNewUser(username, apiKey)
	} else {
		userID, _ := rt.db.GetUserID(username)
		apiKey, _ = rt.db.GetUserKey(userID)
	}

	response := struct {
		Username string `json:"username"`
		APIKey   string `json:"apiKey"`
	}{
		Username: username,
		APIKey:   apiKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
