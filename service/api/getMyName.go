package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
//		"encoding/json"
		"image/gif"
       )

type User struct {
                Id int 
                Name string
                Photo *gif.GIF
                security string
                Conversations []int
}

func (rt *_router) getMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
