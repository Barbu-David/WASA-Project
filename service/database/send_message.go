package database

import "time"

func (db *appdbimpl) SendMessage(senderID int, convID int, textContent string, forwarded bool, timestamp time.Time) error {
	res, err := db.c.Exec(`
		INSERT INTO Messages (sender_id, conv_id, content, forwarded, timestamp)
		VALUES (?, ?, ?, ?, ?)`,
		senderID, convID, textContent, forwarded, timestamp)
	if err != nil {
		return err
	}

	// Retrieve the ID of the newly inserted message.
	msgID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// For every user in the conversation, insert an entry into the SeenList table
	rows, err := db.c.Query(`
		SELECT user_id 
		FROM ConversationMembers 
		WHERE conv_id = ?`, convID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return err
		}

		_, err := db.c.Exec(`
			INSERT INTO SeenList (m_id, user_id, comment, seen, received)
			VALUES (?, ?, ?, ?, ?)`,
			msgID, userID, nil, false, false)
		if err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
