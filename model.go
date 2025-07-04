package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Frames              map[int]string
	Size                *tea.WindowSizeMsg
	WindowSizeChan      chan tea.WindowSizeMsg
	CurrentFrameIndex   int
	LoadingPercentage   int
	FPS                 int
	FAR                 float64
	AsciiConfiguration  AsciiConfig
	ProcessingCompleted bool
	Playing             bool
}

func newModel(
	windowSizeChan chan tea.WindowSizeMsg,
	fps int,
	far float64,
	asciiConfig AsciiConfig,
) model {
	return model{
		CurrentFrameIndex:   0,
		Frames:              make(map[int]string),
		WindowSizeChan:      windowSizeChan,
		FPS:                 fps,
		FAR:                 far,
		AsciiConfiguration:  asciiConfig,
		ProcessingCompleted: false,
		Playing:             false,
	}
}

func (m *model) reset() {
	m.CurrentFrameIndex = 0
	m.Frames = make(map[int]string)
	m.Size = nil
}

func (m *model) loading(p int) {
	m.CurrentFrameIndex = 0
	m.Frames = make(map[int]string)
	m.LoadingPercentage = p
}

func (m *model) setProcessingAsCompleted(oldModel *model) {
	m.ProcessingCompleted = true
	m.Size = oldModel.Size
}

func (m *model) play() {
	m.Playing = true
}

type loadingMsg struct {
	p int
}

type processingCompletedMsg struct {
	m model
}

type playMsg struct{}

type frameMsg struct{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.WindowSizeChan <- msg
		m.Size = &msg
		return m, func() tea.Msg {
			return loadingMsg{
				p: 0,
			}
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeySpace:
			if m.ProcessingCompleted {
				return m, func() tea.Msg {
					return playMsg{}
				}
			}
		case tea.KeyRunes:
			switch msg.String() {
			case "q", "esc", "ctrl+c":
				return m, tea.Quit
			}
		}

	case processingCompletedMsg:
		msg.m.setProcessingAsCompleted(&m)
		return msg.m, nil

	case playMsg:
		if m.Playing {
			return m, nil
		}
		m.play()
		return m, func() tea.Msg {
			return frameMsg{}
		}

	case loadingMsg:
		m.loading(msg.p)
		return m, nil

	case frameMsg:

		if len(m.Frames) == 0 {
			return m, nil
		}

		if m.CurrentFrameIndex < len(m.Frames) {
			m.CurrentFrameIndex++
		} else {
			m.CurrentFrameIndex = 0
		}

		// NOTE: some gifs can have empty frames so we skipped them
		if canShowFrame(m) {
			return m, tick(m.FPS)
		} else {
			return m, func() tea.Msg {
				return frameMsg{}
			}
		}

	}

	return m, nil
}

func canShowFrame(m model) bool {
	return len(m.Frames[m.CurrentFrameIndex]) > 0
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
		if m.ProcessingCompleted {
			res += "Press <space> in order to start playing the gif"
		}
		res = style.Render(res)
	} else if len(m.Frames) == m.CurrentFrameIndex {
		frame := m.Frames[m.CurrentFrameIndex]
		res += frame
	} else {
		frame := m.Frames[m.CurrentFrameIndex]
		res += frame
	}

	return res
}

func tick(fps int) tea.Cmd {
	return tea.Tick(time.Second/time.Duration(fps), func(t time.Time) tea.Msg {
		return frameMsg{}
	})
}
