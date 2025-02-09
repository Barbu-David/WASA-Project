package database

func (db *appdbimpl) DeleteConversationMember(user_id int, conv_id int) error {
	_, err := db.c.Exec(`
		DELETE FROM ConversationMembers
		WHERE conv_id = ? AND user_id = ?`,
		conv_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
