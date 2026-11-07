package main

import "testing"

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
