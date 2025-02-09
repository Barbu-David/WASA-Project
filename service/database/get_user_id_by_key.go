package database

func (db *appdbimpl) GetUserIDbyKey(security_key string) (int, error) {
	var ID int
	err := db.c.QueryRow(`
		SELECT ID FROM Users WHERE security_key = ?`, security_key).Scan(&ID)
	if err != nil {
		return -1, err
	}
	return ID, nil
}
