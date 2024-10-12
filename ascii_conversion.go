package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

type AsciiConfig struct {
	CharacterDensity string
	SetRandomBlank   bool
}

var style = lipgloss.NewStyle()

func imageToAscii(
	img image.Image,
	config AsciiConfig,
) string {
	if img == nil {
		return ""
	}

	ww := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	var wg sync.WaitGroup
	numWorkers := 4
	chunkSize := h / numWorkers
	results := make([]string, numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID, start, end int) {
			defer wg.Done()
			res := ""
			for i := start; i < end; i++ {
				for j := 0; j < ww; j++ {
					co := img.At(j, i)
					rr, gg, bb, _ := co.RGBA()
					r := uint8(rr)
					g := uint8(gg)
					b := uint8(bb)
					hex := rgbToHex(r, g, b)
					c := characterFromRgb(r, g, b, config.CharacterDensity)
					if hex == "#000000" {
						s := style.
							Foreground(lipgloss.Color("#FFFFFF"))
						if config.SetRandomBlank {
							res += s.Render(getRandomToken(config.CharacterDensity))
						} else {
							res += s.Render(string(config.CharacterDensity[0]))
						}
					} else {
						complementaryHex := rgbToHex(255-r, 255-g, 255-b)
						s := style.
							Background(lipgloss.Color(hex)).
							Foreground(lipgloss.Color(complementaryHex))
						res += s.Render(c)
					}
				}
				res = res + "\n"
			}
			results[workerID] = res
		}(w, w*chunkSize, (w+1)*chunkSize)
	}

	wg.Wait()

	finalResult := ""
	for _, result := range results {
		finalResult += result
	}

	return finalResult
}

func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

func rgbToHex(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func characterFromRgb(r, g, b uint8, cd string) string {
	avgF := float64(int(r)+int(g)+int(b)) / 3.0
	avg := uint8(math.Round(avgF))
	len := len(cd)
	i := int(mapValue(avg, 0, 255, 0, uint8(len)))
	if i >= len {
		i = len - 1
	}
	return string(cd[int(i)])
}
