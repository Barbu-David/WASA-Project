package database

func (db *appdbimpl) NewConversation(name string, group bool) (int, error) {
	res, err := db.c.Exec(`
		INSERT INTO Conveversations (name, is_group) 
		VALUES (?, ?)`, name, group)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}
