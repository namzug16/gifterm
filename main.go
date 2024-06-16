package main

import (
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

type frameMsg struct{}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
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
		files, err := m.readFiles()
		if err != nil {
			m.Error = err
			return m, nil
		}

		var fileNames []string

		for _, file := range files {
			if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
				fileNames = append(fileNames, m.Dir+"/"+file.Name())
			}
		}

		sort.Slice(fileNames, func(i, j int) bool {
			return fileNames[i] < fileNames[j]
		})

		m.Files = fileNames

		m.Width = message.(tea.WindowSizeMsg).Width
		m.Height = message.(tea.WindowSizeMsg).Height

		c1 := m.loadImages(m.Files)
		results := make(chan job)

		numWorkers := 30
		var wg sync.WaitGroup

		wg.Add(numWorkers)

		for i := 0; i < numWorkers; i++ {
			go m.worker(i, &wg, c1, results)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		// countInitialFrames := fps * 3
		//
		// wg.Add(countInitialFrames)
		//
		//   fmt.Println("Started ", countInitialFrames)
		//
		// go func() {
		// 	for j := range results {
		//       m.Mx.Lock()
		// 		m.Frames[j.InputPath] = j.Ascii
		//       m.Mx.Unlock()
		// 		wg.Done()
		//       fmt.Println("Done")
		// 	}
		// }()
		//
		// wg.Wait()
		//
		//   fmt.Println("finished waiting ")

		for j := range results {
			m.Frames[j.InputPath] = j.Ascii
		}

		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case frameMsg:
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

	if m.Error != nil {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000FF"))
		return style.Render(m.Error.Error())
	}

	if m.Debug {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
		res = res + style.Render(fmt.Sprintf("Width: %d, Height: %d, Files count: %d Current frame: %d, Frames count: %d\n", m.Width, m.Height, len(m.Files), m.CurrentFrameIndex, len(m.Frames)))
		// NOTE: height needs to be reduce so that debug is shown
	}

	if len(m.Frames) == 0 {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
		res = res + style.Render("Files is fucking empty bro\n")
	} else if len(m.Frames) == m.CurrentFrameIndex {
		frameName := m.Files[m.CurrentFrameIndex-1]

		frame := m.Frames[frameName]

		res += frame
	} else {

		frameName := m.Files[m.CurrentFrameIndex]

		if m.Debug {
			res += frameName + "\n"
		}

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
