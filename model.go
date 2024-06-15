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

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/image/draw"
)

type model struct {
	Error             error
	Frames            map[string]string
	Dir               string
	Files             []string
	Width             int
	Height            int
	CurrentFrameIndex int
	Debug             bool
  Mx *sync.Mutex
}

type frame struct {
	Name string
	Data string
}

func newModel() model {
  var mx sync.Mutex
	return model{
		CurrentFrameIndex: 0,
		Width:             0,
		Height:            0,
		Debug:             true,
		Dir:               "./input",
		Frames:            make(map[string]string),
    Mx: &mx,
	}
}

type job struct {
	InputPath string
	Image     image.Image
	Ascii     string
}

// WARNING: CHANNELS =========================================================
func (m model) loadImages(paths []string) <-chan job {
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

func (m model) resizeImages(input <-chan job) <-chan job {
	out := make(chan job)
	go func() {
		for job := range input {
			job.Image = resizeImage(job.Image, m.Width, m.Height)
			out <- job
		}
		close(out)
	}()
	return out
}

func (m model) imagesToAscii(input <-chan job) <-chan job {
	out := make(chan job)
	go func() {
		for job := range input {
			// job.Ascii = imageToAscii3(job.Image, m.Width, m.Height)
			job.Ascii = imageToAscii3(job.Image)
			out <- job
		}
		close(out)
	}()
	return out
}

// WARNING: =========================================================
func (m model) readFiles() ([]fs.DirEntry, error) {
	files, err := os.ReadDir(m.Dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read input directory. %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("empty fucking directory")
	}
	return files, nil
}

func (m model) worker(
	i int,
	wg *sync.WaitGroup,
	jobs <-chan job,
	result chan<- job,
) {
	defer wg.Done()
	// for j := range jobs {
	// 	j.Image = m.resizeImage(j.Image)
	// 	j.Ascii = m.imageToAscii(j.Image)
	//    result <- j
	// }
	c2 := m.resizeImages(jobs)
	c3 := m.imagesToAscii(c2)
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

func imageToAscii1(img image.Image, ww, h int) string {
	var wg sync.WaitGroup
	numWorkers := 4 // Number of concurrent workers
	chunkSize := h / numWorkers
	results := make([]string, numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID, start, end int) {
			defer wg.Done()
			res := ""
			for i := start; i < end; i++ {
				for j := 0; j < ww; j++ {
					color := img.At(j, i)
					hex := colorToHex(color)
					if hex == "#000000" {
						res = res + " "
					} else {
						style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
						res = res + style.Render("X")
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

func imageToAscii2(img image.Image) string {
	res := ""

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		panic("Not rgba")
	}

	bounds := rgbaImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height

	for i := 0; i < totalPixels; i++ {
		x := i % width
		y := i / width
		offset := y*rgbaImg.Stride + x*4
		r, g, b, _ := rgbaImg.Pix[offset], rgbaImg.Pix[offset+1], rgbaImg.Pix[offset+2], rgbaImg.Pix[offset+3]
		hex := rgbToHex(r, g, b)
		if hex == "#000000" {
			res = res + " "
		} else {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
			res = res + style.Render("X")
		}
	}

	return res
}

func imageToAscii3(img image.Image) string {
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		panic("Not rgba")
	}

	bounds := rgbaImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height

	var mx sync.Mutex
	var wg sync.WaitGroup

	// Divide the total pixels into 4 parts
	numSegments := 6
	segmentSize := totalPixels / numSegments

	r := make(map[int]string, numSegments)

	for i := 0; i < numSegments; i++ {
		start := i * segmentSize
		end := start + segmentSize
		if i == numSegments-1 {
			// Make sure the last segment processes any remaining pixels
			end = totalPixels
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			res := ""
			for i := start; i < end; i++ {
				x := i % width
				y := i / width
				offset := y*rgbaImg.Stride + x*4
				r, g, b, _ := rgbaImg.Pix[offset], rgbaImg.Pix[offset+1], rgbaImg.Pix[offset+2], rgbaImg.Pix[offset+3]
				hex := rgbToHex(r, g, b)
        complementaryHex := rgbToHex(255 - r, 255 - g, 255 -b)
				if hex == "#000000" {
					res += " "
				} else {
					style := lipgloss.
						NewStyle().
						Background(lipgloss.Color(hex)).
						Foreground(lipgloss.Color(complementaryHex))
					c := characterFromRgb(r, r, g)
					res += style.Render(c)
				}
				if x == width-1 {
					res += "\n"
				}
			}
			mx.Lock()
			r[i] = res
			mx.Unlock()
		}(i)
	}

	wg.Wait()

	res := ""
	for i := 0; i < numSegments; i++ {
		res += r[i]
	}

	return res
}

func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

func rgbToHex(r, g, b uint8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

const cDensity = ".:-ix|=+%O#@X"

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
