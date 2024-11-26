package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
		"encoding/json"
       )

func CheckConversationAccess(convId int, userId int) (int, []int, []int) {
	for _, convo := range conversations {
		if convo.Id == convId {
			for _, member := range convo.Members {
				if member == userId {
					return 200, convo.Messages, convo.Members
				}
			}
			return 401, nil, nil
		}
	}
	return 404, nil, nil
}

func FindConvById(convid int) *Conversation {
    for _, foundConv := range conversations {
        if foundConv.Id == convid {
            return &foundConv
        }
    }
    return nil 
}


func DoConvId(tokenStr int) (errorMessage, int){

	var errorm errorMessage 

	foundConversation :=FindConvById(tokenStr);
	if foundConversation==nil {
		errorm.Content="Conversation not found"
			return errorm, 404
	}

	return errorm, 200

}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
	
		json.NewEncoder(w).Encode(details)
		w.WriteHeader(200)
}
