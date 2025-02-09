package database

func (db *appdbimpl) NewConversationMember(user_id int, conv_id int) error {
	_, err := db.c.Exec(`
		INSERT INTO ConversationMembers (conv_id, user_id)
		VALUES (?, ?)`,
		conv_id, user_id)
	if err != nil {
		return err
	}
	return nil
}
