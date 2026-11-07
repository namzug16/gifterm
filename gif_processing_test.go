package main

import (
	"image"
	"image/color"
	"image/gif"
	"testing"
)

func TestComposeFramesAppliesDisposal(t *testing.T) {
	palette := color.Palette{color.Transparent, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255}}

	red := image.NewPaletted(image.Rect(0, 0, 4, 4), palette)
	for i := range red.Pix {
		red.Pix[i] = 1
	}
	green := image.NewPaletted(image.Rect(2, 2, 3, 3), palette)
	green.Pix[0] = 2

	g := &gif.GIF{
		Image:    []*image.Paletted{red, green, green},
		Disposal: []byte{gif.DisposalBackground, gif.DisposalPrevious, gif.DisposalNone},
		Config:   image.Config{Width: 4, Height: 4},
	}

	frames := composeFrames(g)

	checkPixel := func(frame, x, y int, want color.RGBA) {
		t.Helper()
		if got := color.RGBAModel.Convert(frames[frame].At(x, y)); got != want {
			t.Errorf("frame %d pixel (%d,%d) = %v, want %v", frame, x, y, got, want)
		}
	}

	checkPixel(0, 0, 0, color.RGBA{255, 0, 0, 255})
	checkPixel(1, 0, 0, color.RGBA{})
	checkPixel(1, 2, 2, color.RGBA{0, 255, 0, 255})
	checkPixel(2, 2, 2, color.RGBA{0, 255, 0, 255})
}
