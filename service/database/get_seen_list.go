package database

func (db *appdbimpl) GetMessageSeenList(mID int) ([]int, error) {
	rows, err := db.c.Query(`
		SELECT user_id
		FROM SeenList
		WHERE m_id = ? AND seen = ?`, mID, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
