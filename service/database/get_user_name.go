package database

func (db *appdbimpl) GetUserName(userID int) (string, error) {
	var username string
	err := db.c.QueryRow(`
		SELECT username FROM Users WHERE ID = ?`, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
