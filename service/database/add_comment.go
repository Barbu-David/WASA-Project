package database

func (db *appdbimpl) AddComment(sender_id int, m_id int, content string) error {
	_, err := db.c.Exec(`
		UPDATE SeenList
		SET comment = ?
		WHERE m_id = ? AND user_id = ?`, content, m_id, sender_id)
	if err != nil {
		return err
	}
	
	return nil
}
