package database

import "image/gif"

func (db *appdbimpl) SetConversationPhoto(conv_id int, photo *gif.GIF) error {

	photo_bytes, err := encodeGIF(photo)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`
			UPDATE Conversations 
			SET gif_photo = ? 
			WHERE ID = ?`, photo_bytes, conv_id)
	if err != nil {
		return err
	}
	return nil
}
