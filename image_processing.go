package main

import (
	"image"

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

func getNewImageBounds(srcW, srcH, w, h int, car float64) (int, int) {
	var newW, newH int

	rw := float64(srcW) / float64(w)
	rh := float64(srcH) / float64(h)

	if rw > rh {
		newH = int(float64(w) * float64(srcH) / float64(srcW) * car)
		newW = w
	} else {
		newW = int(float64(h) * float64(srcW) / float64(srcH) * car)
		newH = h
	}

	return newW, newH
}
