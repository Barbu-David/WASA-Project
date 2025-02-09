package database

func (db *appdbimpl) GetConversationName(conv_id int) (string, error) {
	var name string
	err := db.c.QueryRow(`
		SELECT name FROM Conversations WHERE ID = ?`, conv_id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
