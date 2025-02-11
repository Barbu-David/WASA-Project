package database

func (db *appdbimpl) DeleteMessage(m_id int) error {
	_, err := db.c.Exec(`
		DELETE FROM Messages
		WHERE ID = ?`,
		m_id)
	if err != nil {
		return err
	}
	_, err = db.c.Exec(`
		DELETE FROM SeenList
		WHERE m_id = ?`,
		m_id)
	if err != nil {
		return err
	}
	return nil
}
