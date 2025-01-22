package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.wrap(rt.doLogin))


	rt.router.GET("/users", rt.wrap(rt.getMaxId))
	rt.router.PUT("/users/:id/name", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:id/name", rt.wrap(rt.getMyUserName))
	rt.router.PUT("/users/:id/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users/:id/photo", rt.wrap(rt.getMyPhoto))

	rt.router.GET("/conversations/:convid", rt.wrap(rt.getConversation))
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))

	rt.router.POST("/conversations/:convid/", rt.wrap(rt.sendMessage))
	rt.router.DELETE("/conversations/:convid/messages/:messageId", rt.wrap(rt.deleteMessage))
	rt.router.POST("/conversations/:convid/messages/:messageId", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/conversations/:convid/messages/:messageId/comments", rt.wrap(rt.uncommentMessage))
	rt.router.POST("/conversations/:convid/messages/:messageId/comments", rt.wrap(rt.commentMessage))
	rt.router.GET("/conversations/:convid/messages/:messageId", rt.wrap(rt.getMessage))
	rt.router.POST("/conversations/:convid/messages/:messageId/photo", rt.wrap(rt.getMessagePhoto))

	rt.router.PUT("/conversations", rt.wrap(rt.startConversation))
	rt.router.PUT("/conversations/:convid/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/conversations/:convid/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/conversations/:convid/name", rt.wrap(rt.getGroupName))
	rt.router.GET("/conversations/:convid/photo", rt.wrap(rt.getGroupPhoto))
	rt.router.PUT("/conversations/:convid/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/conversations/:convid/members", rt.wrap(rt.leaveGroup))

	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
