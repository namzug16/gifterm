package main

import (
	"image"
	"image/color"
	"testing"
)

func TestGetNewImageBounds(t *testing.T) {
	tests := []struct {
		name         string
		srcW, srcH   int
		w, h         int
		car          float64
		fit          ObjectFit
		wantW, wantH int
	}{
		{"contain very wide fills width", 400, 100, 80, 24, 2, ObjectFitContain, 80, 10},
		{"contain wide fills width", 200, 100, 80, 24, 2, ObjectFitContain, 80, 20},
		{"contain square fills height", 100, 100, 80, 24, 2, ObjectFitContain, 48, 24},
		{"contain tall fills height", 100, 400, 80, 24, 2, ObjectFitContain, 12, 24},
		{"contain square cells", 100, 100, 80, 24, 1, ObjectFitContain, 24, 24},
		{"cover very wide overflows width", 400, 100, 80, 24, 2, ObjectFitCover, 192, 24},
		{"cover wide overflows width", 200, 100, 80, 24, 2, ObjectFitCover, 96, 24},
		{"cover square overflows height", 100, 100, 80, 24, 2, ObjectFitCover, 80, 40},
		{"cover tall overflows height", 100, 400, 80, 24, 2, ObjectFitCover, 80, 160},
		{"cover square cells", 100, 100, 80, 24, 1, ObjectFitCover, 80, 80},
		{"fill stretches wide", 400, 100, 80, 24, 2, ObjectFitFill, 80, 24},
		{"fill stretches tall", 100, 400, 80, 24, 2, ObjectFitFill, 80, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotH := getNewImageBounds(tt.srcW, tt.srcH, tt.w, tt.h, tt.car, tt.fit)
			if gotW != tt.wantW || gotH != tt.wantH {
				t.Errorf("got %dx%d, want %dx%d", gotW, gotH, tt.wantW, tt.wantH)
			}
		})
	}
}

func TestResizeImage(t *testing.T) {
	white := color.RGBA{255, 255, 255, 255}
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := range 100 {
		for x := range 100 {
			src.Set(x, y, white)
		}
	}

	tests := []struct {
		fit           ObjectFit
		wantBlackBars bool
	}{
		{ObjectFitContain, true},
		{ObjectFitCover, false},
		{ObjectFitFill, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.fit), func(t *testing.T) {
			dst := resizeImage(src, 80, 24, 2, tt.fit)

			if got := dst.Bounds(); got != image.Rect(0, 0, 80, 24) {
				t.Fatalf("bounds = %v, want 80x24", got)
			}

			// a square gif on an 80x24 terminal with car 2 only fills 48 columns when contained
			_, _, _, edgeAlpha := dst.At(0, 12).RGBA()
			if gotBlackBars := edgeAlpha == 0; gotBlackBars != tt.wantBlackBars {
				t.Errorf("black bars = %v, want %v", gotBlackBars, tt.wantBlackBars)
			}
			if got := color.RGBAModel.Convert(dst.At(40, 12)); got != white {
				t.Errorf("center = %v, want white", got)
			}
		})
	}
}

func TestParseObjectFit(t *testing.T) {
	for _, value := range []string{"contain", "cover", "fill"} {
		if fit, err := parseObjectFit(value); err != nil || string(fit) != value {
			t.Errorf("parseObjectFit(%q) = %q, %v", value, fit, err)
		}
	}
	if _, err := parseObjectFit("scale-down"); err == nil {
		t.Error("parseObjectFit(\"scale-down\") should fail")
	}
}
