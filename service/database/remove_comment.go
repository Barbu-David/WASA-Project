package database

func (db *appdbimpl) RemoveComment(sender_id int, m_id int) error {
	_, err := db.c.Exec(`
		UPDATE SeenList
		SET comment = NULL
		WHERE m_id = ? AND user_id = ?`, m_id, sender_id)
	if err != nil {
		return err
	}
	
	return nil
}
