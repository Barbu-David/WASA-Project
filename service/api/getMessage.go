package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
	"wasatext/service/api/reqcontext"
)

func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authHeader := r.Header.Get("Authorization")

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization format"})
		return
	}

	token := authHeader[len(bearerPrefix):]

	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Empty token"})
		return
	}

	conv_id_param := ps.ByName("convid")

	conv_id, err := strconv.Atoi(conv_id_param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid conversation id"})
		return
	}

	// User must exist

	user_id, err := rt.db.GetUserIDbyKey(token)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Auth error"})
		return
	}

	// And be a member of the conversation

	member, err := rt.db.IsMemberConversation(user_id, conv_id)

	if err != nil || member != true {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Auth error"})
		return
	}

	m_id_param := ps.ByName("messageId")

	m_id, err := strconv.Atoi(m_id_param)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid mid"})
		return
	}

	err = rt.db.SeeMessage(user_id, m_id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Database failed")
	}

	mark1, err := rt.db.IsDeliveredToAll(m_id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Database failed")
	}

	mark2, err := rt.db.IsSeenByAll(m_id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Database failed")
	}

	var mark string

	if mark2 != false {
		mark = "✓✓"
	} else if mark1 != false {
		mark = "✓"
	} else {
		mark = "..."
	}

	sender_id, content, fwded, stamp, photoContent, err := rt.db.GetMessage(m_id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Database failed")
	}

	ownersList, commentList, err := rt.db.GetMessageCommentList(m_id)

	if err != nil {
		ctx.Logger.WithError(err).Error("Database failed")
	}

	response := struct {
		StringContent string    `json:"stringContent"`
		SenderId      int       `json:"senderId"`
		Timestamp     time.Time `json:"timestamp"`
		Checkmark     string    `json:"checkmark"`
		Forwarded     bool      `json:"forwarded"`
		PhotoContent  bool      `json:"photoContent"`
		Comments      []string  `json:"comments"`
		CommentOwners []int     `json:"comment_owners"`
	}{
		StringContent: content,
		SenderId:      sender_id,
		Timestamp:     stamp,
		Checkmark:     mark,
		Forwarded:     fwded,
		PhotoContent:  photoContent,
		Comments:      commentList,
		CommentOwners: ownersList,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Error encoding response"})
		ctx.Logger.WithError(err).Error("Encoding failed")
	}
}
