package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.doLogin)
	
	rt.router.PUT("/users/:id/name", rt.setMyUserName)
	rt.router.PUT("/users/:id/photo", rt.setMyPhoto)

	rt.router.GET("/conversations/:convid", rt.getConversation)
	rt.router.GET("/conversations", rt.getMyConversations)

	rt.router.POST("/conversations/:convid/", rt.sendMessage)
	rt.router.DELETE("/conversations/:convid/messages/:messageId", rt.deleteMessage)
        rt.router.POST("/conversations/:convid/messages/:messageId", rt.forwardMessage)
        rt.router.DELETE("/conversations/:convid/messages/:messageId/comments", rt.uncommentMessage)
        rt.router.POST("/conversations/:convid/messages/:messageId/comments", rt.commentMessage)

	rt.router.PUT("/conversations/:convid/name", rt.setGroupName)
        rt.router.PUT("/conversations/:convid/photo", rt.setGroupPhoto)
	rt.router.PUT("/conversations/:convid/members", rt.addToGroup)
        rt.router.DELETE("/conversations/:convid/members", rt.leaveGroup)

	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
