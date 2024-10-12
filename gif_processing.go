package main

import (
	"fmt"
	"image"
	"image/gif"
	"os"
)

func readGif(path string) ([]*image.Paletted, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening GIF:", err)
		return nil, err
	}
	defer file.Close()

	gifImage, err := gif.DecodeAll(file)
	if err != nil {
		fmt.Println("Error decoding GIF:", err)
		return nil, err
	}

	return gifImage.Image, nil
}
