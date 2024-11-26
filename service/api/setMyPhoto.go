package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
		"encoding/json"
		"image/gif"
       )

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "application/json")

	tokenStr := r.Header.Get("security_token")

	user, errMessage, statusCode := DoSecurity(tokenStr)
	if user==nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errMessage)
		return 	
	}
	var newPhoto *gif.GIF

	err := json.NewDecoder(r.Body).Decode(&newPhoto)
	if err != nil {
		errMessage.Content="Invalid request body"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	user.Photo=newPhoto

	w.WriteHeader(204)
}
