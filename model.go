package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Frames              []string
	Size                *tea.WindowSizeMsg
	RenderRequests      chan<- renderRequest
	RenderGeneration    int
	CurrentFrameIndex   int
	LoadingPercentage   int
	FPS                 int
	CAR                 float64
	AsciiConfiguration  AsciiConfig
	ProcessingCompleted bool
	Playing             bool
}

func newModel(
	renderRequests chan<- renderRequest,
	fps int,
	car float64,
	asciiConfig AsciiConfig,
) model {
	return model{
		RenderRequests:     renderRequests,
		FPS:                fps,
		CAR:                car,
		AsciiConfiguration: asciiConfig,
	}
}

// Messages carry the RenderGeneration they belong to, so the ones from
// before the latest resize are ignored.

type loadingMsg struct {
	Generation int
	Percentage int
}

type processingCompletedMsg struct {
	Generation int
	Frames     []string
}

type frameMsg struct {
	Generation int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.RenderGeneration++
		m.RenderRequests <- renderRequest{Size: msg, Generation: m.RenderGeneration}
		m.Size = &msg
		m.Frames = nil
		m.CurrentFrameIndex = 0
		m.LoadingPercentage = 0
		m.ProcessingCompleted = false
		m.Playing = false

	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			if m.ProcessingCompleted && !m.Playing {
				m.Playing = true
				return m, tick(m.FPS, m.RenderGeneration)
			}
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case loadingMsg:
		if msg.Generation == m.RenderGeneration {
			m.LoadingPercentage = msg.Percentage
		}

	case processingCompletedMsg:
		if msg.Generation == m.RenderGeneration {
			m.Frames = msg.Frames
			m.ProcessingCompleted = true
		}

	case frameMsg:
		if msg.Generation == m.RenderGeneration {
			m.CurrentFrameIndex = (m.CurrentFrameIndex + 1) % len(m.Frames)
			return m, tick(m.FPS, m.RenderGeneration)
		}
	}

	return m, nil
}

func (m model) View() string {
	res := ""

	if !m.Playing {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
		if !m.ProcessingCompleted {
			res += "Loading your frames... "
			res += fmt.Sprintf("%d", m.LoadingPercentage)
			res += "%\n"
		} else {
			res += "Frames Completed\n"
		}
		if m.Size != nil {
			res += fmt.Sprintf("Width: %d, Height: %d\n", m.Size.Width, m.Size.Height)
		}
		res += "Character Density: " + m.AsciiConfiguration.CharacterDensity + "\n"
		res += "Random Blank: " + fmt.Sprint(m.AsciiConfiguration.SetRandomBlank) + "\n"
		res += "FPS: " + fmt.Sprint(m.FPS) + "\n"
		res += "Cell Aspect Ratio: " + fmt.Sprint(m.CAR) + "\n"
		if m.ProcessingCompleted {
			res += "Press <space> in order to start playing the gif"
		}
		res = style.Render(res)
	} else {
		res += m.Frames[m.CurrentFrameIndex]
	}

	return res
}

func tick(fps int, generation int) tea.Cmd {
	return tea.Tick(time.Second/time.Duration(fps), func(t time.Time) tea.Msg {
		return frameMsg{Generation: generation}
	})
}
