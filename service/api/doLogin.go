package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
		"encoding/json"
		"crypto/rand"
		"math/big"
       )

func FindUserByName(name string) *User {
	for _, user := range users {
		if user.Name == name {
			return &user
		}
	}
	return nil 
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString() string {

length := 16
		result := make([]byte, length)
		for i := range result {
			num, err:= rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
				result[i] = charset[num.Int64()]
				if(err!=nil){ return string(result)
				}
		}
	return string(result)
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "application/json")

	var newName Username
	var errMessage errorMessage
	err := json.NewDecoder(r.Body).Decode(&newName)
	if err != nil {
		errMessage.Content="Invalid request body"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	if len(newName.Content) < 3 || len(newName.Content) > 16 {
		errMessage.Content="Invalid lenght"
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(errMessage)
		return		
	}

	user := FindUserByName(newName.Content)
	var lg Login
	if user==nil {
	lg.SecToken=GenerateRandomString()
	lg.Id=len(users)+1
	}
	if(user!=nil) {
	lg.SecToken=user.security
	lg.Id=user.Id
	}
	json.NewEncoder(w).Encode(lg)
	w.WriteHeader(201)
}
