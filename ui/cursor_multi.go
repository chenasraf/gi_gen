package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type CursorM[T any, S any] struct {
	question string
	choices  []*Choice[T]
	cursor   int
	selected S
	done     bool
	quitting bool
}

func (m *CursorM[T, S]) Select(i int) tea.Cmd {
	return nil
}

func (m *CursorM[T, S]) Init() tea.Cmd {
	return nil
}

func (m *CursorM[T, S]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case " ":
			m.Select(m.cursor)
			// _, exists := m.selected[m.cursor]
			// if exists {
			// 	delete(m.selected, m.cursor)
			// } else {
			// 	m.selected[m.cursor] = struct{}{}
			// }
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func isExists[T any, S any](m *CursorM[T, S]) func(i int) bool {
	return func(i int) bool {
		// _, exists := m.selected[i]
		// return exists
		return false
	}
}

func (m *CursorM[T, S]) View() string {
	if m.done || m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString(m.question + "\n")
	// s.WriteString(fmt.Sprintf("%d selected\n", len(m.selected)))
	s.WriteString(viewRowsWindow(ViewWindow[T]{
		height:     10,
		maxWidth:   60,
		offset:     5,
		isSelected: isExists[T](m),
		cursor:     m.cursor,
		choices:    m.choices,
		checkbox:   true,
		border:     true,
	}))
	s.WriteString("up/down/j/k - move cursor, space - select, enter - confirm\n")

	return s.String()
}
