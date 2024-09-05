package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type CursorM[T any] struct {
	question     string
	choices      []*Choice[T]
	cursor       int
	selected     map[int]struct{}
	done         bool
	quitting     bool
	lastMovement int
}

func (m CursorM[T]) Init() tea.Cmd {
	return nil
}

func (m CursorM[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.choices) - 1
			}
			m.lastMovement = -1
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
			m.lastMovement = 1
		case " ":
			_, exists := m.selected[m.cursor]
			if exists {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func isExists[T any](m CursorM[T]) func(i int) bool {
	return func(i int) bool {
		_, exists := m.selected[i]
		return exists
	}
}

func (m CursorM[T]) View() string {
	if m.done || m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString(m.question + "\n")
	s.WriteString(fmt.Sprintf("%d selected\n", len(m.selected)))
	s.WriteString(viewRowsWindow(ViewWindow[T]{
		height:       10,
		maxWidth:     60,
		offset:       5,
		isSelected:   isExists[T](m),
		cursor:       m.cursor,
		lastMovement: m.lastMovement,
		choices:      m.choices,
		checkbox:     true,
		border:       true,
	}))
	s.WriteString("up/down/j/k - move cursor, space - select, enter - confirm\n")

	return s.String()
}
