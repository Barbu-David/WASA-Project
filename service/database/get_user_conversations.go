package database

func (db *appdbimpl) GetUserConversations(userID int) ([]int, error) {
	query := `SELECT conv_id FROM ConversationMembers WHERE user_id = ?;`

	rows, err := db.c.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convIDs []int
	for rows.Next() {
		var convID int
		if err := rows.Scan(&convID); err != nil {
			return nil, err
		}
		convIDs = append(convIDs, convID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return convIDs, nil
}
