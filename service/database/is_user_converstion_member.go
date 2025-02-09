package database

func (db *appdbimpl) IsMemberConversation(user_id int, conv_id int) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(1) 
		FROM ConversationMembers 
		WHERE user_id = ? AND conv_id = ?`, user_id, conv_id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
