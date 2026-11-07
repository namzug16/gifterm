package main

import (
	"image"
	"image/gif"
	"os"
)

func readGif(path string) ([]*image.Paletted, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gifImage, err := gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}

	return gifImage.Image, nil
}
