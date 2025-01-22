package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	//              "encoding/json"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
}
