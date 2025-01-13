package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
//		"encoding/json"
		"image/gif"
       )

type Conversation struct {
                Id int 
                Name string
                Group bool
                Photo *gif.GIF
                Members []int
                Messages []int
}

type ConvDetails struct {
                Conversations []int
                Members []int
}

/*
func GetUserConversationIDs(userId int) []int {
	return userConversationIDs
}
*/

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
