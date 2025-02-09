package database

func (db *appdbimpl) GetConversationUsers(conv_id int) ([]int, error) {
	query := `SELECT user_id FROM ConversationMembers WHERE conv_id = ?;`

	rows, err := db.c.Query(query, conv_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}
