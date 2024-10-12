package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	Frames            map[int]string
	Size              *tea.WindowSizeMsg
	WindowSizeChan    chan tea.WindowSizeMsg
	CurrentFrameIndex int
	LoadingPercentage int
	FPS               int
	FAR               float64
}

func newModel(
	windowSizeChan chan tea.WindowSizeMsg,
	fps int,
	far float64,
) model {
	return model{
		CurrentFrameIndex: 0,
		Frames:            make(map[int]string),
		WindowSizeChan:    windowSizeChan,
		FPS:               fps,
		FAR:               far,
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

type loadingMsg struct {
	p int
}

type playMsg struct {
	m model
}

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
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case playMsg:
		return msg.m, func() tea.Msg {
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

func tick(fps int) tea.Cmd {
	return tea.Tick(time.Second/time.Duration(fps), func(t time.Time) tea.Msg {
		return frameMsg{}
	})
}
