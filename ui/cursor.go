package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Cursor[T any] struct {
	question     string
	choices      []*Choice[T]
	cursor       int
	selected     int
	lastMovement int
}

func (m Cursor[T]) Init() tea.Cmd {
	return nil
}

func (m Cursor[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return nil, tea.Quit
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
		case "enter", " ":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}

	return m, nil
}

func isSelected[T any](m Cursor[T]) func(i int) bool {
	return func(i int) bool {
		return m.selected == i
	}
}

func (m Cursor[T]) View() string {
	var s strings.Builder
	s.WriteString(m.question + "\n\n")
	s.WriteString("\nup/down/j/k - move cursor, space/enter - confirm\n")
	s.WriteString(viewRowsWindow(ViewWindow[T]{
		height:       10,
		offset:       5,
		isSelected:   isSelected[T](m),
		cursor:       m.cursor,
		lastMovement: m.lastMovement,
		choices:      m.choices,
		checkbox:     false,
	}))

	return s.String()
}
