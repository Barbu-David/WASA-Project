package database

import "time"

func (db *appdbimpl) SendMessage(senderID int, convID int, textContent string, forwarded bool, timestamp time.Time) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec(`
            INSERT INTO Messages (sender_id, conv_id, content, forwarded, timestamp)
            VALUES (?, ?, ?, ?, ?)`,
		senderID, convID, textContent, forwarded, timestamp)
	if err != nil {
		tx.Rollback()
		return err
	}

	msgID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := tx.Query(`
            SELECT user_id 
            FROM ConversationMembers 
            WHERE conv_id = ?`, convID)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			tx.Rollback()
			return err
		}

		_, err := tx.Exec(`
                INSERT INTO SeenList (m_id, user_id, comment, seen, received)
                VALUES (?, ?, ?, ?, ?)`,
			msgID, userID, nil, false, false)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
