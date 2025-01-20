package database

func (db *appdbimpl) GetUserKey(userID int) (string, error) {
	var key string
	err := db.c.QueryRow(`
		SELECT security_key FROM Users WHERE ID = ?`, userID).Scan(&key)
	if err != nil {
		return "", err
	}
	return key, nil
}
