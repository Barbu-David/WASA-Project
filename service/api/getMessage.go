package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
//		"encoding/json"
		"image/gif"
		"time"
       )

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


func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
}
