package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.doLogin)

	rt.router.PUT("/users/:id/name", rt.setMyUserName)
	rt.router.GET("/users/:id/name", rt.getMyUserName)
	rt.router.PUT("/users/:id/photo", rt.setMyPhoto)
	rt.router.GET("/users/:id/photo", rt.getMyPhoto)

	rt.router.GET("/conversations/:convid", rt.getConversation)
	rt.router.GET("/conversations", rt.getMyConversations)

	rt.router.POST("/conversations/:convid/", rt.sendMessage)
	rt.router.DELETE("/conversations/:convid/messages/:messageId", rt.deleteMessage)
	rt.router.POST("/conversations/:convid/messages/:messageId", rt.forwardMessage)
	rt.router.DELETE("/conversations/:convid/messages/:messageId/comments", rt.uncommentMessage)
	rt.router.POST("/conversations/:convid/messages/:messageId/comments", rt.commentMessage)
	rt.router.GET("/conversations/:convid/messages/:messageId", rt.getMessage)
	rt.router.POST("/conversations/:convid/messages/:messageId/photo", rt.getMessagePhoto)

	rt.router.PUT("/conversations/:convid/name", rt.setGroupName)
	rt.router.PUT("/conversations/:convid/photo", rt.setGroupPhoto)
	rt.router.GET("/conversations/:convid/name", rt.getGroupName)
	rt.router.GET("/conversations/:convid/photo", rt.getGroupPhoto)
	rt.router.PUT("/conversations/:convid/members", rt.addToGroup)
	rt.router.DELETE("/conversations/:convid/members", rt.leaveGroup)

	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
