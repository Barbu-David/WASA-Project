package api

import (
		"github.com/julienschmidt/httprouter"
		"net/http"
//		"encoding/json"
//		"crypto/rand"
//		"math/big"
       )

type Login struct {
        Id int
        SecToken string

}


func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

