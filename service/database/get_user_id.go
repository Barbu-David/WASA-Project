package database

func (db *appdbimpl) GetUserID(username string) (int, error) {
	var ID int
	err := db.c.QueryRow(`
		SELECT ID FROM Users WHERE username = ?`, username).Scan(&ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}
