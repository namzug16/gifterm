package main

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

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
		// NOTE: read files
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

    //FIX: make this a worker xdxdxdxd

		c1 := m.loadImages(m.Files)
		c2 := m.resizeImages(c1)
		c3 := m.imagesToAscii(c2)

		for j := range c3 {
			m.Frames[j.InputPath] = j.Ascii
		}

		return m, tick()
		// return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tickMsg:
		// NOTE: next frame from generated ones
		if m.CurrentFrameIndex < len(m.Files) {
			m.CurrentFrameIndex++
			return m, tick()
		} else {
			return m, nil
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
	return tea.Tick(time.Second/20, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
