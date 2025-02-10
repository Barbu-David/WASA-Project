package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {

	rt.router.POST("/session", rt.wrap(rt.doLogin)) // *

	rt.router.GET("/users", rt.wrap(rt.getMaxId))               // *
	rt.router.PUT("/users/:id/name", rt.wrap(rt.setMyUserName)) // *
	rt.router.GET("/users/:id/name", rt.wrap(rt.getMyUserName)) // *
	rt.router.PUT("/users/:id/photo", rt.wrap(rt.setMyPhoto))   // MARTI
	rt.router.GET("/users/:id/photo", rt.wrap(rt.getMyPhoto))   // MARTI

	rt.router.PUT("/new_conversation", rt.wrap(rt.startConversation)) // *

	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))      // *
	rt.router.GET("/conversations/:convid", rt.wrap(rt.getConversation)) // *

	rt.router.POST("/conversations/:convid/", rt.wrap(rt.sendMessage))  // LUNI
	rt.router.DELETE("/conversations/:convid/messages/:messageId", rt.wrap(rt.deleteMessage)) // LUNI
	rt.router.POST("/conversations/:convid/messages/:messageId", rt.wrap(rt.forwardMessage)) //LUNI
	rt.router.DELETE("/conversations/:convid/messages/:messageId/comments", rt.wrap(rt.uncommentMessage)) // LUNI
	rt.router.POST("/conversations/:convid/messages/:messageId/comments", rt.wrap(rt.commentMessage)) // LUNI
	rt.router.GET("/conversations/:convid/messages/:messageId", rt.wrap(rt.getMessage)) // LUNI
	rt.router.GET("/conversations/:convid/messages/:messageId/photo", rt.wrap(rt.getMessagePhoto)) // MARTI

	rt.router.PUT("/conversations/:convid/name", rt.wrap(rt.setGroupName))     // *
	rt.router.PUT("/conversations/:convid/photo", rt.wrap(rt.setGroupPhoto))   // MARTI
	rt.router.GET("/conversations/:convid/name", rt.wrap(rt.getGroupName))     // *
	rt.router.GET("/conversations/:convid/photo", rt.wrap(rt.getGroupPhoto))   // MARTI
	rt.router.PUT("/conversations/:convid/members", rt.wrap(rt.addToGroup))    // *
	rt.router.DELETE("/conversations/:convid/members", rt.wrap(rt.leaveGroup)) // *

	rt.router.GET("/liveness", rt.liveness) // *

	return rt.router
}
