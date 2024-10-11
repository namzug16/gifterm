package main

import (
	"image"
	"image/png"
	"os"
	"testing"

	"golang.org/x/image/draw"
)

func TestProcessGif(t *testing.T) {
	gifPath := "output.gif"

	imgs, err := readGif(gifPath)
	if err != nil {
		t.Log("===== err")
		t.Log(err)
		return
	}

	img := imgs[0]

	srcBounds := img.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()
	t.Log("Bound of image", srcBounds)
	t.Log("Bound of image", srcW, srcH)

	w := 178

	h := 48

	dst := image.Rect(0, 0, w, h)

	t.Log("Bound of view", dst)
	t.Log("Bound of view", w, h)

	newW, newH := getNewImageBounds(srcW, srcH, w, h)

	offsetX := (w - newW) / 2
	offsetY := (h - newH) / 2

	dstRect := image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH)

	t.Log("New Bounds of image inside view", dstRect)
	t.Log("New Bounds of image inside view", newW, newH)

	dstImage := image.NewRGBA(dst)

	draw.NearestNeighbor.Scale(dstImage, dstRect, img, img.Bounds(), draw.Over, nil)

	t.Log("Image bounds", dstImage.Bounds(), dstImage.Bounds().Max)

	outputFile, err := os.Create("output.png")
	if err != nil {
		t.Log("Failed to create output file:", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, dstImage)
	if err != nil {
		t.Log("Failed to save image:", err)
		return
	}

	t.Log("Image saved as output.png successfully")
}
