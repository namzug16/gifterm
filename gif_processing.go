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

	return composeFrames(gifImage), nil
}

// composeFrames draws each frame onto the full GIF canvas, since frames only
// store what changed, and applies each frame's disposal before the next one.
func composeFrames(g *gif.GIF) []image.Image {
	canvas := image.NewRGBA(image.Rect(0, 0, g.Config.Width, g.Config.Height))
	composed := make([]image.Image, len(g.Image))

	for i, frame := range g.Image {
		var canvasBeforeFrame *image.RGBA
		if g.Disposal[i] == gif.DisposalPrevious {
			canvasBeforeFrame = cloneRGBA(canvas)
		}

		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		composed[i] = cloneRGBA(canvas)

		switch g.Disposal[i] {
		case gif.DisposalBackground:
			draw.Draw(canvas, frame.Bounds(), image.Transparent, image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			canvas = canvasBeforeFrame
		}
	}

	return composed
}

func cloneRGBA(img *image.RGBA) *image.RGBA {
	clone := image.NewRGBA(img.Rect)
	copy(clone.Pix, img.Pix)
	return clone
}
