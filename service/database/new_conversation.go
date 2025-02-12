package database

func (db *appdbimpl) NewConversation(name string, group bool) (int, error) {

	defaultPhoto := createDefaultGIF()

	bytes, err := encodeGIF(defaultPhoto)
	if err != nil {
		return 0, err
	}

	res, err := db.c.Exec(`
		INSERT INTO Conversations (name, is_group, gif_photo) 
		VALUES (?, ?, ?)`, name, group, bytes)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}
