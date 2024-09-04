package prompt

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type CursorM[T any] struct {
	question string
	choices  []*Choice[T]
	cursor   int
	selected map[int]struct{}
}

func (m CursorM[T]) Init() tea.Cmd {
	return nil
}

func (m CursorM[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case " ":
			_, exists := m.selected[m.cursor]
			if exists {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m CursorM[T]) View() string {
	s := m.question + "\n\n"

	// max := m.cursor + 5
	h := 10
	offset := 2
	for r := range h {
		i := m.cursor - offset + r
		// fmt.Printf("r: %v, i: %v\n", r, i)
		if i < 0 || i > len(m.choices)-1 {
			s += "\n"
			continue
		}
		choice := m.choices[i]
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, exists := m.selected[i]; exists {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.Label)
	}

	s += "\nup/down/j/k - move cursor, space - select, enter - confirm\n"

	return s
}
