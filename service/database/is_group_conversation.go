package database

func (db *appdbimpl) IsGroupConversation(conv_id int) (bool, error) {
	var group bool
	err := db.c.QueryRow(`
		SELECT is_group FROM Conversations WHERE ID = ?`, conv_id).Scan(&group)
	if err != nil {
		return false, err
	}
	return group, nil
}
