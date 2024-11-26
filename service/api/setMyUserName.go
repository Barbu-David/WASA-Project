package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
		"encoding/json"
		"image/gif"
		"time"
       )

type User struct {
		Id int
		Name string
		Photo *gif.GIF
		security string
		Conversations []int
}

type ConvDetails struct {
		Conversations []int
		Members []int
}


type Conversation struct {
		Id int
		Name string
		Group bool
		Photo *gif.GIF
		Members []int
		Messages []int
}

type Message struct {
		Id int
		ConvId int
		Content_type bool
		Content_text string
		Content_photo *gif.GIF
		Timestamp time.Time
		Sender int
		Seen []int
		Forwarded bool
		Comment string
}

type Login struct {
	Id int
	SecToken string

}

type errorMessage struct {
	Content string
}

type Username struct {
	Content string
}


var users  []User
var conversations  []Conversation
var messages [][]Message

func FindUserByToken(tokenStr string) *User {
    for _, user := range users {
        if user.security == tokenStr {
            return &user
        }
    }
    return nil 
}

func FindUserIndexByToken(tokenStr string) int {
    for index, user := range users {
        if user.security == tokenStr {
            return index
        }
    }
    return -1
}



func DoSecurity(tokenStr string) (*User, errorMessage, int){

	var errorm errorMessage	
		if tokenStr == "" {
			errorm.Content="Empty security token"
				return nil, errorm, 400
		}

	foundUser :=FindUserByToken(tokenStr);
	if foundUser==nil {
		errorm.Content="User not found"
			return nil, errorm, 404	
	}
	
	return foundUser, errorm, 0

}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "application/json")

	tokenStr := r.Header.Get("security_token")

	user, errMessage, statusCode := DoSecurity(tokenStr)
	if user==nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errMessage)
		return 	
	}
	var newName Username

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

	user.Name=newName.Content

	w.WriteHeader(204)
}
