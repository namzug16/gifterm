package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps = 24
)

type loadingMsg struct {
	p int
}

type playMsg struct {
	m model
}

type updateFramesMsg struct {
	frames map[string]string
	files  []string
}

type frameMsg struct{}

func main() {
	dir := "input/"

	windowSizeChan := make(chan tea.WindowSizeMsg)
	ctx, cancel := context.WithCancel(context.Background())

	m := newModel(
		windowSizeChan,
	)

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func(m *model) {
		files, _ := readFiles(dir)

		var fileNames []string

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
				fileNames = append(fileNames, dir+file.Name())
			}
		}

		sort.Slice(fileNames, func(i, j int) bool {
			return fileNames[i] < fileNames[j]
		})

		for size := range windowSizeChan {
			cancel()
			ctx, cancel = context.WithCancel(context.Background())

			go func(size tea.WindowSizeMsg) {
				select {
				case <-ctx.Done():
					return
				default:
					c1 := loadImages(fileNames)
					results := make(chan job)

					numWorkers := 10
					var wg sync.WaitGroup

					wg.Add(numWorkers)

					for i := 0; i < numWorkers; i++ {
						go worker(i, &wg, c1, results, size.Width, size.Height)
					}

					go func() {
						wg.Wait()
						close(results)
					}()

					go func() {
						frames := make(map[string]string)

						for j := range results {
							frames[j.InputPath] = j.Ascii
							pe := int((float32(len(frames)) / float32(len(fileNames))) * 100)
							p.Send(loadingMsg{
								p: pe,
							})
						}

						m.Files = fileNames
						m.Frames = frames

						p.Send(playMsg{
							m: *m,
						})
					}()
				}
			}(size)
		}
	}(&m)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.WindowSizeChan <- msg
		return m.reset(), func() tea.Msg {
			return loadingMsg{
				p: 0,
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case playMsg:
		return msg.m, func() tea.Msg {
			return frameMsg{}
		}

	case loadingMsg:
		m.LoadingPercentage = msg.p
    m.Files = make([]string, 0)
		return m, nil

	case updateFramesMsg:
		m.Frames = msg.frames
		m.Files = msg.files
		return m, nil

	case frameMsg:
    if len(m.Files) == 0 {
      return m, nil
    }

		if m.CurrentFrameIndex < len(m.Files) {
			m.CurrentFrameIndex++
			return m, tick()
		} else {
			m.CurrentFrameIndex = 0
			return m, tick()
			// return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	res := ""

	if len(m.Frames) == 0 {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
		res += "Loading your frames... "
		res += fmt.Sprintf("%d", m.LoadingPercentage)
		res += "%"
		res = style.Render(res)
	} else if len(m.Frames) == m.CurrentFrameIndex {
		frameName := m.Files[m.CurrentFrameIndex-1]

		frame := m.Frames[frameName]

		res += frame
	} else {

		frameName := m.Files[m.CurrentFrameIndex]

		frame := m.Frames[frameName]

		if len(frame) == 0 {
			res += "LOADING FRAMES"
		} else {
			res += frame
		}
	}

	return res
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return frameMsg{}
	})
}

// countInitialFrames := fps * 5
//
// if countInitialFrames > len(fileNames) {
// 	countInitialFrames = len(fileNames)
// }
// go func() {
// 	c := countInitialFrames
// 	isFirst := true
// 	frames := make(map[string]string)
// 	for j := range results {
// 		frames[j.InputPath] = j.Ascii
// 		if isFirst {
// 			if c > 0 {
// 				c--
// 				pe := int((float32(len(frames)) / float32(countInitialFrames)) * 100)
// 				p.Send(loadingMsg{
// 					p: pe,
// 				})
// 				if c == 0 && countInitialFrames == len(fileNames) {
// 					m.Files = fileNames
// 					m.Frames = frames
// 					p.Send(playMsg{
// 						m: *m,
// 					})
// 				}
// 			} else {
// 				isFirst = false
// 				c = fps
// 				m.Files = fileNames
// 				m.Frames = frames
// 				p.Send(playMsg{
// 					m: *m,
// 				})
// 			}
// 		} else {
// 			if c > 0 {
// 				c--
// 			} else {
// 				p.Send(updateFramesMsg{
// 					frames: frames,
// 					files:  fileNames,
// 				})
// 				c = fps
// 			}
// 		}
// 	}
// }()
