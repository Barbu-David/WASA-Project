package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
		"encoding/json"
       )

func GetUserConversationIDs(userId int) []int {
	var userConversationIDs []int

	for _, convo := range conversations {
		for _, member := range convo.Members {
			if member == userId {
				userConversationIDs = append(userConversationIDs, convo.Id)
				break
			}
		}
	}

	return userConversationIDs
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "application/json")

	tokenStr := r.Header.Get("security_token")

	user, errMessage, statusCode := DoSecurity(tokenStr)
	if user==nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errMessage)
		return 	
	}

	userConvos := GetUserConversationIDs(user.Id)
	json.NewEncoder(w).Encode(userConvos)

	w.WriteHeader(200)
}
