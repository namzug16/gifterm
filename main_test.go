package main

import (
	"testing"
	"time"
)

func IndexCharacter(t *testing.T) {
	r := uint8(200)
	g := uint8(120)
	b := uint8(200)

	c := characterFromRgb(r, g, b)

	t.Log("CHARACTER : ", c)
}

func ProcessingSingleImageAscii(t *testing.T) {
	img := readImage("input/frame_0001.png")
	start := time.Now()
	imageToAscii(img)
	stop := time.Since(start)
	t.Log("Image ascii: ", stop)
}

func TestProcessMultipleImages(t *testing.T) {
}
