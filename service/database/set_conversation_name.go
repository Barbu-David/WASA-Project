package database

func (db *appdbimpl) SetConversationName(conv_id int, name string) error {
	_, err := db.c.Exec(`
			UPDATE Conversations 
			SET name = ? 
			WHERE ID = ?`, name, conv_id)
	if err != nil {
		return err
	}
	return nil
}
