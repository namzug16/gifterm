package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fps = 12
)

type loadingMsg struct {
	p int
}

type playMsg struct {
	m model
}

type frameMsg struct{}

func main() {
	gifPath := "output.gif"

	windowSizeChan := make(chan tea.WindowSizeMsg)
	// ctx, cancel := context.WithCancel(context.Background())

	m := newModel(
		windowSizeChan,
	)

	p := tea.NewProgram(m, tea.WithAltScreen())

	go func(m *model) {
		images, err := readGif(gifPath)
		// FIX: handle error
		if err != nil {
			return
		}

		size := <-windowSizeChan

		go func(size tea.WindowSizeMsg) {
			c1 := chanFromImages(images)
			results := make(chan job)

			numWorkers := 10
			var wg sync.WaitGroup

			wg.Add(numWorkers)

			for i := 0; i < numWorkers; i++ {
				go worker(&wg, c1, results, size.Width, size.Height)
			}

			go func() {
				wg.Wait()
				close(results)
			}()

			go func() {
				frames := make(map[int]string)

				for j := range results {
					frames[j.Index] = j.Ascii
					pe := int((float32(len(frames)) / float32(len(images))) * 100)
					p.Send(loadingMsg{
						p: pe,
					})
				}

				m.Frames = frames

				p.Send(playMsg{
					m: *m,
				})
			}()
		}(size)
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
		m.Size = &msg
		return m, func() tea.Msg {
			return frameMsg{}
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
		return m.loading(msg.p), nil

	case frameMsg:
		if len(m.Frames) == 0 {
			return m, nil
		}

    //FIX: why am I doing - 1 ???
		if m.CurrentFrameIndex < len(m.Frames)-1 {
			m.CurrentFrameIndex++
			return m, tick()
		} else {
			m.CurrentFrameIndex = 0
			return m, tick()
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
		if m.Size != nil {
			res += fmt.Sprintf(", Width: %d, Height: %d", m.Size.Width, m.Size.Height)
		}
		res = style.Render(res)
	} else if len(m.Frames) == m.CurrentFrameIndex {
		frame := m.Frames[m.CurrentFrameIndex]
		res += frame
	} else {
		frame := m.Frames[m.CurrentFrameIndex]
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
