package database

func (db *appdbimpl) AddNewUser(username string, securityKey string) (int, error) {
	res, err := db.c.Exec(`
		INSERT INTO Users (username, security_key) 
		VALUES (?, ?)`, username, securityKey)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

