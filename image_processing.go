package main

import (
	"fmt"
	"image"
	"math"

	"golang.org/x/image/draw"
)

// ObjectFit uses the CSS object-fit vocabulary for how the gif fills the terminal.
type ObjectFit string

const (
	ObjectFitContain ObjectFit = "contain"
	ObjectFitCover   ObjectFit = "cover"
	ObjectFitFill    ObjectFit = "fill"
)

func parseObjectFit(value string) (ObjectFit, error) {
	switch fit := ObjectFit(value); fit {
	case ObjectFitContain, ObjectFitCover, ObjectFitFill:
		return fit, nil
	}
	return "", fmt.Errorf("unknown fit %q, must be one of: contain, cover, fill", value)
}

func resizeImage(img image.Image, w, h int, car float64, fit ObjectFit) image.Image {
	srcBounds := img.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()

	newW, newH := getNewImageBounds(srcW, srcH, w, h, car, fit)

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	// offsets are negative when the image overflows (cover), Scale clips it to dst
	offsetX := (w - newW) / 2
	offsetY := (h - newH) / 2

	dstRect := image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH)

	draw.CatmullRom.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)

	return dst
}

// getNewImageBounds scales the image for w columns by h rows, where each cell is car times taller than wide.
func getNewImageBounds(srcW, srcH, w, h int, car float64, fit ObjectFit) (int, int) {
	if fit == ObjectFitFill {
		return w, h
	}

	newW := w
	newH := int(math.Round(float64(w) * float64(srcH) / float64(srcW) / car))

	// contain shrinks to the height when the image is too tall, cover grows to it when too short
	if (fit == ObjectFitContain && newH > h) || (fit == ObjectFitCover && newH < h) {
		newH = h
		newW = int(math.Round(float64(h) * float64(srcW) / float64(srcH) * car))
	}

	return newW, newH
}
