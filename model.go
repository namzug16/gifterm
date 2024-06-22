package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"math"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/image/draw"
)

type model struct {
	Frames            map[string]string
	Files             []string
	CurrentFrameIndex int
	LoadingPercentage int
	WindowSizeChan    chan tea.WindowSizeMsg
}

type frame struct {
	Name string
	Data string
}

func newModel(
	windowSizeChan chan tea.WindowSizeMsg,
) model {
	return model{
		CurrentFrameIndex: 0,
		Frames:            make(map[string]string),
		WindowSizeChan:    windowSizeChan,
	}
}

func (m model) reset() model {
	return model{
		CurrentFrameIndex: 0,
		Frames:            make(map[string]string),
		WindowSizeChan:    m.WindowSizeChan,
	}
}

type job struct {
	InputPath string
	Image     image.Image
	Ascii     string
}

// WARNING: CHANNELS =========================================================
func loadImages(paths []string) <-chan job {
	out := make(chan job)
	go func() {
		for _, p := range paths {
			job := job{
				InputPath: p,
				Image:     readImage(p),
			}
			out <- job
		}
		close(out)
	}()
	return out
}

func resizeImages(input <-chan job, w, h int) <-chan job {
	out := make(chan job)
	go func() {
		for job := range input {
			job.Image = resizeImage(job.Image, w, h)
			out <- job
		}
		close(out)
	}()
	return out
}

func imagesToAscii(input <-chan job) <-chan job {
	out := make(chan job)
	go func() {
		for job := range input {
			job.Ascii = imageToAscii(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

// WARNING: =========================================================
func readFiles(dir string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read input directory. %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("empty fucking directory")
	}
	return files, nil
}

func worker(
	i int,
	wg *sync.WaitGroup,
	jobs <-chan job,
	result chan<- job,
	w,
	h int,
) {
	defer wg.Done()
	c2 := resizeImages(jobs, w, h)
	c3 := imagesToAscii(c2)
	for j := range c3 {
		result <- j
	}
}

func readImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("CANT READ SHIT")
		fmt.Println(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("CANT decode SHIT")
		fmt.Println(path)
		fmt.Println(err)
	}
	return img
}

func resizeImage(img image.Image, w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

var style = lipgloss.NewStyle()

func imageToAscii(img image.Image) string {
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
					c := characterFromRgb(r, g, b)
					if hex == "#000000" {
						s := style.
							Foreground(lipgloss.Color("#FFFFFF"))
						res += s.Render(string(cDensity[0]))
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

const cDensity = ".,:-=i|%O#@$X"

func characterFromRgb(r, g, b uint8) string {
	avgF := float64(int(r)+int(g)+int(b)) / 3.0
	avg := uint8(math.Round(avgF))
	len := len(cDensity)
	i := int(mapValue(avg, 0, 255, 0, uint8(len)))
	if i >= len {
		i = len - 1
	}
	return string(cDensity[int(i)])
}

func mapValue(
	value uint8,
	minIn uint8,
	maxIn uint8,
	minOut uint8,
	maxOut uint8,
) uint8 {
	finalValue := value

	if value > maxIn {
		finalValue = maxIn
	} else if value < minIn {
		finalValue = minIn
	}

	a := float64(maxIn-finalValue) / float64(maxIn-minIn)

	b := a * float64(maxOut-minOut)

	c := math.Round(b)

	if c < 0 {
		c = 0
	} else if c > 255 {
		c = 255
	}

	return maxOut - uint8(c)
}
