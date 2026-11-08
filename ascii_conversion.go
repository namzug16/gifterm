package main

import (
	"fmt"
	"image"
	"math/rand/v2"
	"strings"

	"charm.land/lipgloss/v2"
)

type AsciiConfig struct {
	CharacterDensity string
	SetRandomBlank   bool
	OnlyForeground   bool
}

var style = lipgloss.NewStyle()

var blankStyle = style.Foreground(lipgloss.Color("#FFFFFF"))

func imageToAscii(
	img image.Image,
	config AsciiConfig,
) string {
	characters := []rune(config.CharacterDensity)
	bounds := img.Bounds()

	var res strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if y > bounds.Min.Y {
			res.WriteByte('\n')
		}
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rr, gg, bb, _ := img.At(x, y).RGBA()
			r, g, b := uint8(rr>>8), uint8(gg>>8), uint8(bb>>8)

			if r == 0 && g == 0 && b == 0 {
				blank := characters[0]
				if config.SetRandomBlank {
					blank = characters[rand.IntN(len(characters))]
				}
				res.WriteString(blankStyle.Render(string(blank)))
				continue
			}

			var s lipgloss.Style
			if config.OnlyForeground {
				s = style.Foreground(lipgloss.Color(rgbToHex(r, g, b)))
			} else {
				s = style.
					Background(lipgloss.Color(rgbToHex(r, g, b))).
					Foreground(lipgloss.Color(rgbToHex(255-r, 255-g, 255-b)))
			}
			res.WriteString(s.Render(characterFromRgb(r, g, b, characters)))
		}
	}

	return res.String()
}

func rgbToHex(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// characterFromRgb picks a character by brightness, splitting 0-255 into equal buckets.
func characterFromRgb(r, g, b uint8, characters []rune) string {
	brightness := (int(r) + int(g) + int(b)) / 3
	return string(characters[brightness*len(characters)/256])
}
