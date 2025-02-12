package database

import "image/gif"

func (db *appdbimpl) SetUserPhoto(userID int, photo *gif.GIF) error {

	photo_bytes, err := encodeGIF(photo)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`
			UPDATE Users 
			SET gif_photo = ? 
			WHERE ID = ?`, photo_bytes, userID)

	if err != nil {
		return err
	}

	return nil
}
