package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"os"

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
}

type frame struct {
	Name string
	Data string
}

func newModel() model {
	return model{
		CurrentFrameIndex: 0,
		Width:             0,
		Height:            0,
		Debug:             true,
		Dir:               "./input",
		Frames:            make(map[string]string),
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
				Image:     m.readImage(p),
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
			job.Image = m.resizeImage(job.Image)
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
			job.Ascii = m.imageToAscii(job.Image)
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
	jobs <-chan job,
	result chan<- job,
) {
	for j := range jobs {
		func() {

			result <- frame{
				Name: j,
				Data: data,
			}
		}()
	}
}

func (m model) readImage(path string) image.Image {
	file, err := os.Open(path)
  if err != nil {
    fmt.Println("CANT READ SHIT");
    fmt.Println(err);
  }
	defer file.Close()
	img, err := png.Decode(file)
  if err != nil {
    fmt.Println("CANT decode SHIT");
    fmt.Println(path);
    fmt.Println(err);
  }
	return img
}

func (m model) resizeImage(img image.Image) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func (m model) imageToAscii(img image.Image) string {
	res := ""
	height := m.Height

	if m.Debug {
		height = height - 2
	}

	for i := 0; i < height; i++ {
		for j := 0; j < m.Width; j++ {
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

	return res
}

func colorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}



	// case tea.WindowSizeMsg:
	// 	// NOTE: read files
	// 	files, err := m.readFiles()
	// 	if err != nil {
	// 		m.Error = err
	// 		return m, nil
	// 	}
	//
	// 	m.Files = files
	//
	// 	const numWorkers = 10
	// 	jobs := make(chan string, len(files))
	// 	results := make(chan frame, len(files))
	// 	errChan := make(chan error, 1)
	//
	// 	m.Width = message.(tea.WindowSizeMsg).Width
	// 	m.Height = message.(tea.WindowSizeMsg).Height
	//
	// 	var wg sync.WaitGroup
	//
	// 	wg.Add(numWorkers)
	//
	// 	for w := 0; w < numWorkers; w++ {
	// 		go func() {
	// 			defer wg.Done()
	// 			m.worker(jobs, results, errChan)
	// 		}()
	// 	}
	//
	// 	go func() {
	// 		for i := 0; i < len(files); i++ {
	// 			jobs <- filepath.Join(m.Dir, files[i].Name())
	// 		}
	// 		close(jobs)
	// 	}()
	//
	// 	go func() {
	// 		wg.Wait()
	// 		close(results)
	// 	}()
	//
	// 	var mx sync.Mutex
	//
	// 	wg.Add(len(files))
	//
	// 	go func() {
	// 		for r := range results {
	// 			// fmt.Println("ADDED")
	// 			// fmt.Println(r.Name)
	// 			// fmt.Println(len(r.Data))
	// 			mx.Lock()
	// 			m.Frames[r.Name] = r.Data
	// 			mx.Unlock()
	// 			wg.Done()
	// 		}
	// 	}()
	//
	// 	wg.Wait()
	//
	// 	fmt.Println("FINISHED")
	// 	return m, tick()
