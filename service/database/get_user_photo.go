package database

import "image/gif"

func (db *appdbimpl) GetUserPhoto(userID int) (*gif.GIF, error) {
	var photo *gif.GIF
	var photo_bytes []byte

	err := db.c.QueryRow(`
		SELECT gif_photo FROM Users WHERE ID = ?`, userID).Scan(&photo_bytes)
	if err != nil {
		return nil, err
	}

	photo, err = decodeGIF(photo_bytes)

	if err != nil {
		return nil, err
	}

	return photo, nil
}
