package database

import "image/gif"

func (db *appdbimpl) SetUserPhoto(userID int, photo *gif.GIF) error {
	_, err := db.c.Exec(`
			UPDATE Users 
			SET gif_photo = ? 
			WHERE ID = ?`, photo, userID)
	if err != nil {
		return err
	}
	return nil
}
