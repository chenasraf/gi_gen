package prompt

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Cursor[T any] struct {
	question string
	choices  []*Choice[T]
	cursor   int
	selected int
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
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter", " ":
			m.selected = m.cursor
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Cursor[T]) View() string {
	s := m.question + "\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice.Label)
	}

	s += "\nup/down/j/k - move cursor, space/enter - confirm\n"

	return s
}
