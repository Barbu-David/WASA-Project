package database

import "time"

func (db *appdbimpl) GetMessage(m_id int) (int, string, bool, time.Time, bool, error) {
	var sender_id int
	var content string
	var forwarded bool
	var is_photo bool
	var timestamp time.Time

	err := db.c.QueryRow(`
		SELECT sender_id, content, forwarded, timestamp, is_photo FROM Messages WHERE ID = ?`, m_id).Scan(&sender_id, &content, &forwarded, &timestamp, &is_photo)
	if err != nil {
		return -1, "", false, time.Time{}, false, err
	}
	return sender_id, content, forwarded, timestamp, is_photo, nil
}
