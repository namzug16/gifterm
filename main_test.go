package main

import (
	"testing"
	"time"
)

func TestProcessGif(t *testing.T) {
	start := time.Now()

	gifPath := "output.gif"

	imgs, err := readGif(gifPath)

  if (err != nil) {
    t.Log("Shit is not wokring")
    t.Log(err)
    return
  }

	stop := time.Since(start)
	start = time.Now()
	t.Log("Load gif in: ", stop)

	resizedImgs := resizeGifImgs(imgs, 100, 50)
	stop = time.Since(start)
	start = time.Now()
	t.Log("Resized imgs in: ", stop)

	imgsToAscii(resizedImgs)

	stop = time.Since(start)
	t.Log("Ascii in: ", stop)

}

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
