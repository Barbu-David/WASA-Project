package database

import "time"

func (db *appdbimpl) GetMessage(m_id int) (int, string, bool, time.Time, error) {
	var sender_id int
	var content string
	var forwarded bool
	var timestamp time.Time
	err := db.c.QueryRow(`
		SELECT sender_id, content, forwarded, timestamp  FROM Messages WHERE ID = ?`, m_id).Scan(&sender_id, &content, &forwarded, &timestamp)
	if err != nil {
		return -1, "", false, time.Time{}, err
	}
	return sender_id, content, forwarded, timestamp, nil
}
