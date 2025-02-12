package database

import "image/gif"

func (db *appdbimpl) SetConversationPhoto(conv_id int, photo *gif.GIF) error {
	_, err := db.c.Exec(`
			UPDATE Conversations 
			SET gif_photo = ? 
			WHERE ID = ?`, photo, conv_id)
	if err != nil {
		return err
	}
	return nil
}
