package database

import "time"
import "image/gif"

func (db *appdbimpl) SendMessage(senderID int, convID int, textContent string, forwarded bool, timestamp time.Time, is_photo bool, photo *gif.GIF) error {
	tx, err := db.c.Begin()
	if err != nil {
		return err
	}

	var photo_bytes []byte

	if is_photo {
		photo_bytes, err = encodeGIF(photo)
		if err != nil {
			return err
		}
	}

	res, err := tx.Exec(`
            INSERT INTO Messages (sender_id, conv_id, content, forwarded, timestamp, is_photo, gif_photo)
            VALUES (?, ?, ?, ?, ?, ?, ?)`,
		senderID, convID, textContent, forwarded, timestamp, is_photo, photo_bytes)
	if err != nil {
		return err
	}

	msgID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	rows, err := tx.Query(`
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

		_, err := tx.Exec(`
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

	return tx.Commit()
}
