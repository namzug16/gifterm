package main

import (
	"image"

	"golang.org/x/image/draw"
)

func resizeImage(img image.Image, w, h int, far float64) image.Image {
	srcBounds := img.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()

	newW, newH := getNewImageBounds(srcW, srcH, w, h, far)

	dst := image.NewRGBA(image.Rect(0, 0, w, h))

	offsetX := (w - newW) / 2
	offsetY := (h - newH) / 2

	dstRect := image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH)

	draw.CatmullRom.Scale(dst, dstRect, img, img.Bounds(), draw.Over, nil)

	return dst
}

func getNewImageBounds(srcW, srcH, w, h int, far float64) (int, int) {
	var newW, newH int

	rw := float64(srcW) / float64(w)
	rh := float64(srcH) / float64(h)

  // NOTE: far represents the font aspect ratio, this value cannot be read from the terminal so it requires manual settings :/ 
	if rw > rh {
		newH = int(float64(w) * float64(srcH) / float64(srcW) * far)
		newW = w
	} else {
		newW = int(float64(h) * float64(srcW) / float64(srcH) * far)
		newH = h
	}

	return newW, newH
}
