package database

func (db *appdbimpl) ReceiveMessage(user_id int, m_id int) error {
	_, err := db.c.Exec(`
		UPDATE SeenList
		SET received = ?
		WHERE m_id = ? AND user_id = ?`, true, m_id, user_id)
	if err != nil {
		return err
	}

	return nil
}
