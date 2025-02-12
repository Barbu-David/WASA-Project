package database

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
)

func createDefaultGIF() *gif.GIF {
	palette := color.Palette{
		color.RGBA{0, 0, 255, 255},
	}
	img := image.NewPaletted(image.Rect(0, 0, 100, 100), palette)
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.SetColorIndex(x, y, 0)
		}
	}

	return &gif.GIF{
		Image:           []*image.Paletted{img},
		Delay:           []int{100},
		LoopCount:       0,
		BackgroundIndex: 0,
	}
}

func encodeGIF(g *gif.GIF) ([]byte, error) {
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, g); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeGIF(data []byte) (*gif.GIF, error) {
	r := bytes.NewReader(data)
	return gif.DecodeAll(r)
}

func (db *appdbimpl) AddNewUser(username string, securityKey string) (int, error) {

	defaultPhoto := createDefaultGIF()

	bytes, err := encodeGIF(defaultPhoto)
	if err != nil {
		return 0, err
	}

	res, err := db.c.Exec(`
		INSERT INTO Users (username, security_key, gif_photo) 
		VALUES (?, ?, ?)`, username, securityKey, bytes)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}
