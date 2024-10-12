package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"

	"golang.org/x/image/draw"
)

func ProcessGif(t *testing.T) {
	gifPath := "output.gif"

	imgs, err := readGif(gifPath)
	if err != nil {
		t.Log("===== err")
		t.Log(err)
		return
	}

	for i := 0; i < len(imgs); i++ {
		img := imgs[i]
		srcBounds := img.Bounds()
		srcW, srcH := srcBounds.Dx(), srcBounds.Dy()
		t.Log("Bound of image", srcBounds)
		t.Log("Bound of image", srcW, srcH)

		fileName := fmt.Sprintf("output_original_%d.png", i)

		// NOTE: STORE GIF IMAGE
		outputFile, err := os.Create(fileName)
		if err != nil {
			t.Log("Failed to create output file:", err)
			return
		}
		defer outputFile.Close()

		err = png.Encode(outputFile, img)
		if err != nil {
			t.Log("Failed to save image:", err)
			return
		}

		t.Log("Image saved successfully", fileName)
		// NOTE:

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

		fileName = fmt.Sprintf("output_%d.png", i)

		outputFile, err = os.Create(fileName)
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

		t.Log("Image saved successfully", fileName)
	}
}

func ConvertPalettedToRGBA(img *image.Paletted) *image.RGBA {
	bounds := img.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)
	return rgbaImg
}

func TestImageToASCII(t *testing.T) {
  	path := "output_7.png"
    file, err := os.Open(path) 
    if err != nil {
        t.Log("Error opening the image:", err)
        return
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        t.Log("Error decoding the image:", err)
        return
    }

  ascii := imageToAscii(img)

			t.Log(ascii)

}
