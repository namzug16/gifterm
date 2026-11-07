package main

import (
	"image"
	"image/gif"
	"os"

	"golang.org/x/image/draw"
)

func readGif(path string) ([]image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gifImage, err := gif.DecodeAll(file)
	if err != nil {
		return nil, err
	}

	return composeFrames(gifImage.Image), nil
}

// composeFrames draws each frame over the previous one, since GIF frames only store what changed.
func composeFrames(frames []*image.Paletted) []image.Image {
	bounds := frames[0].Bounds()
	composed := make([]image.Image, len(frames))
	previous := image.NewRGBA(bounds)

	for i, frame := range frames {
		current := image.NewRGBA(bounds)
		draw.Draw(current, bounds, previous, bounds.Min, draw.Src)
		draw.Draw(current, bounds, frame, bounds.Min, draw.Over)
		composed[i] = current
		previous = current
	}

	return composed
}
