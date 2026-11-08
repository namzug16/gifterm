package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestModel() model {
	return newModel(make(chan renderRequest, 1), 10, 2, ObjectFitContain, AsciiConfig{CharacterDensity: " .#"})
}

func TestSpaceStartsPlaybackOnlyAfterProcessing(t *testing.T) {
	space := tea.KeyPressMsg{Code: tea.KeySpace, Text: " "}

	m := newTestModel()
	updated, cmd := m.Update(space)
	if updated.(model).Playing || cmd != nil {
		t.Fatal("space must not start playback before frames are processed")
	}

	m.ProcessingCompleted = true
	m.Frames = []string{"a", "b"}
	updated, cmd = m.Update(space)
	if !updated.(model).Playing {
		t.Fatal("space should start playback once frames are processed")
	}
	if cmd == nil {
		t.Fatal("starting playback should schedule a tick")
	}
}

func TestQuitKeys(t *testing.T) {
	keys := map[string]tea.KeyPressMsg{
		"q":      {Code: 'q', Text: "q"},
		"esc":    {Code: tea.KeyEscape},
		"ctrl+c": {Code: 'c', Mod: tea.ModCtrl},
	}

	for name, key := range keys {
		t.Run(name, func(t *testing.T) {
			_, cmd := newTestModel().Update(key)
			if cmd == nil {
				t.Fatal("expected a quit command")
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Fatalf("expected tea.QuitMsg, got %T", cmd())
			}
		})
	}
}

func TestWindowSizeRequestsRender(t *testing.T) {
	requests := make(chan renderRequest, 1)
	m := newModel(requests, 10, 2, ObjectFitContain, AsciiConfig{})
	m.Playing = true

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	request := <-requests
	if request.Size.Width != 80 || request.Size.Height != 24 || request.Generation != 1 {
		t.Fatalf("unexpected render request %+v", request)
	}
	if updated.(model).Playing {
		t.Fatal("resize should stop playback until the new frames are ready")
	}
}

func TestViewUsesAltScreen(t *testing.T) {
	view := newTestModel().View()
	if !view.AltScreen {
		t.Fatal("view should request the alt screen")
	}
	if !strings.Contains(view.Content, "Loading your frames") {
		t.Fatalf("unexpected view content %q", view.Content)
	}
}
