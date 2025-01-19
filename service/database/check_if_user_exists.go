package database

func (db *appdbimpl) CheckIfUserExists(username string) (bool, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(1) 
		FROM Users 
		WHERE username = ?`, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
