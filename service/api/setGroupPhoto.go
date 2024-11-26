package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
		"encoding/json"
		"image/gif"
       )

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		w.Header().Set("content-type", "application/json")

		tokenStr := r.Header.Get("security_token")

		user, errMessage, statusCode := DoSecurity(tokenStr)
		if user==nil {
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(errMessage)
			return 	
		}

		var details ConvDetails 
		var convid int
		var status int	
		
		err := json.NewDecoder(r.Body).Decode(&convid)
	
		if err != nil {
			errMessage.Content="Invalid request body"
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		errMessage, status = DoConvId(convid)	
		if status==404 {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(errMessage)
			return		
		}


		status, details.Conversations, details.Members = CheckConversationAccess(convid, user.Id)

		if status==401 {
			errMessage.Content="Not part of the conversation"
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		if status==404 {
			errMessage.Content="Conversation not found"
				w.WriteHeader(404)
				json.NewEncoder(w).Encode(errMessage)
				return
		}
		var newPhoto *gif.GIF

		err = json.NewDecoder(r.Body).Decode(&newPhoto)
		if err != nil {
			errMessage.Content="Invalid request body"
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(errMessage)
			return
		}

		if conversations[convid].Group==false {
			errMessage.Content="Not a Group"
			w.WriteHeader(403)
			json.NewEncoder(w).Encode(errMessage)
			return
		}
		conversations[convid].Photo=newPhoto
		w.WriteHeader(204)
}
