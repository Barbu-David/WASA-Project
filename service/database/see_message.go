package database

func (db *appdbimpl) SeeMessage(userID, mID int) error {
	_, err := db.c.Exec(`
		UPDATE SeenList
		SET seen = ?
		WHERE m_id = ? AND user_id = ?`, true, mID, userID)
	if err != nil {
		return err
	}
	return nil
}
