package main

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

func resizeImage(img image.Image, w, h int, car float64) image.Image {
	srcBounds := img.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()

	newW, newH := getNewImageBounds(srcW, srcH, w, h, car)

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	offsetX := (w - newW) / 2
	offsetY := (h - newH) / 2

	dstRect := image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH)

	draw.CatmullRom.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)

	return dst
}

// getNewImageBounds fits the image into w columns by h rows, where each cell is car times taller than wide.
func getNewImageBounds(srcW, srcH, w, h int, car float64) (int, int) {
	newW := w
	newH := int(math.Round(float64(w) * float64(srcH) / float64(srcW) / car))

	if newH > h {
		newH = h
		newW = int(math.Round(float64(h) * float64(srcW) / float64(srcH) * car))
	}

	return newW, newH
}
