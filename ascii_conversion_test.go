package main

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestImageToAsciiEmitsTrueColor(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	got := imageToAscii(img, AsciiConfig{CharacterDensity: " .#"})

	// Background is the pixel color, foreground its inverse.
	if !strings.Contains(got, "48;2;255;0;0") || !strings.Contains(got, "38;2;0;255;255") {
		t.Fatalf("expected truecolor escapes, got %q", got)
	}
}

func BenchmarkImageToAscii(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 60))
	for y := range 60 {
		for x := range 200 {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y * 4), B: uint8(x + y), A: 255})
		}
	}
	config := AsciiConfig{CharacterDensity: " .:-=+*#%@"}

	for b.Loop() {
		imageToAscii(img, config)
	}
}

func TestCharacterFromRgb(t *testing.T) {
	characters := []rune(" ░▒▓█")

	tests := []struct {
		name    string
		r, g, b uint8
		want    string
	}{
		{"black is first", 0, 0, 0, " "},
		{"white is last", 255, 255, 255, "█"},
		{"mid grey is middle", 128, 128, 128, "▒"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := characterFromRgb(tt.r, tt.g, tt.b, characters); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
