package database

func (db *appdbimpl) GetUserKey(userID int) (string, error) {
	var key string
	err := db.c.QueryRow(`
		SELECT username FROM Users WHERE ID = ?`,  key).Scan(&key)
	if err != nil {
		return "", err  
		}
	return key, nil  
}

