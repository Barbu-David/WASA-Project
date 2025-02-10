package database

func (db *appdbimpl) GetMessageCommentList(mID int) ([]int, []string, error) {
	rows, err := db.c.Query(`
		SELECT user_id, comment
		FROM SeenList
		WHERE m_id = ? AND comment IS NOT NULL`, mID)
	if err != nil {
		return nil, nil,  err
	}
	defer rows.Close()

	var userIDs []int
	var comments []string
	for rows.Next() {
		var userID int
		var comment string
		if err := rows.Scan(&userID, &comment); err != nil {
			return nil, nil, err
		}
		userIDs = append(userIDs, userID)
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return userIDs, comments, nil
}
