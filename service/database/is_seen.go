package database

func (db *appdbimpl) IsSeenByAll(m_id int) (bool, error) {
	query := `SELECT COUNT(*) FROM SeenList WHERE m_id = ? AND seen = FALSE`

	var count int
	err := db.c.QueryRow(query, m_id).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
