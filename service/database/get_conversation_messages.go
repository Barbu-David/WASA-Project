package database

func (db *appdbimpl) GetConversationMessages(conv_id int) ([]int, error) {
	query := `SELECT ID FROM Messages WHERE conv_id = ?;`

	rows, err := db.c.Query(query, conv_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mIDs []int
	for rows.Next() {
		var mID int
		if err := rows.Scan(&mID); err != nil {
			return nil, err
		}
		mIDs = append(mIDs, mID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mIDs, nil
}
