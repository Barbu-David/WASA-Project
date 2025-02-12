package database

import "time"

func (db *appdbimpl) GetMessageLatest(conv_id int) (string, time.Time, bool, error) {
	var content string
	var timestamp time.Time
	var isp bool

	err := db.c.QueryRow(`
		SELECT content, timestamp, is_photo FROM Messages WHERE ID = (SELECT MAX(ID) FROM Messages where conv_id = ?)`, conv_id).Scan(&content, &timestamp, &isp)
	if err != nil {
		return "", time.Time{}, false, err
	}
	return content, timestamp, isp, nil
}
