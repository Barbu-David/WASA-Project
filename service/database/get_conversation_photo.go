package database

import "image/gif"

func (db *appdbimpl) GetConversationPhoto(conv_id int) (*gif.GIF, error) {
	var photo_bytes []byte
	var photo *gif.GIF
	err := db.c.QueryRow(`
		SELECT gif_photo FROM Conversations WHERE ID = ?`, conv_id).Scan(&photo_bytes)
	if err != nil {
		return nil, err
	}

	photo, err = decodeGIF(photo_bytes)

	if err != nil {
		return nil, err
	}

	return photo, nil
}
