package database

func (db *appdbimpl) GetMaxUserID() (int, error) {
	query := "SELECT IFNULL(MAX(id), 0) FROM Users"
	var maxID int

	err := db.c.QueryRow(query).Scan(&maxID)
	if err != nil {
		return 0, err
	}

	return maxID, nil
}
