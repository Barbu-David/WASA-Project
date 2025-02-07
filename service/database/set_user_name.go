package database

func (db *appdbimpl) SetUserName(userID int, username string) error {
	_, err := db.c.Exec(`
			UPDATE Users 
			SET username = ? 
			WHERE ID = ?`, username, userID)
	if err != nil {
		return err
	}
	return nil
}
