package main

import "testing"

func TestGetNewImageBounds(t *testing.T) {
	tests := []struct {
		name         string
		srcW, srcH   int
		w, h         int
		car          float64
		wantW, wantH int
	}{
		{"very wide fills width", 400, 100, 80, 24, 2, 80, 10},
		{"wide fills width", 200, 100, 80, 24, 2, 80, 20},
		{"square fills height", 100, 100, 80, 24, 2, 48, 24},
		{"tall fills height", 100, 400, 80, 24, 2, 12, 24},
		{"square cells", 100, 100, 80, 24, 1, 24, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotH := getNewImageBounds(tt.srcW, tt.srcH, tt.w, tt.h, tt.car)
			if gotW != tt.wantW || gotH != tt.wantH {
				t.Errorf("got %dx%d, want %dx%d", gotW, gotH, tt.wantW, tt.wantH)
			}
		})
	}
}
