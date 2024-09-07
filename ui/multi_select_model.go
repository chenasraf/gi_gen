package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type MultiSelectModel[T any] struct {
	items    []*Choice[T]
	selected map[int]struct{}
}

func (m *MultiSelectModel[T]) Items() []*Choice[T] {
	return m.items
}

func (m *MultiSelectModel[T]) Select(i int) {
	if _, exists := m.selected[i]; exists {
		delete(m.selected, i)
	} else {
		m.selected[i] = struct{}{}
	}
}

func (m *MultiSelectModel[T]) GetSelectedCount() int {
	return len(m.selected)
}

func (m *MultiSelectModel[T]) IsSelected(i int) bool {
	_, exists := m.selected[i]
	return exists
}

func (m *MultiSelectModel[T]) Style() ListStyle {
	style := DefaultListStyle()
	style.DrawCheckbox = true
	style.DrawBorder = true
	return style
}

func (m *MultiSelectModel[T]) Update(msg tea.Msg) (ListModel[T], tea.Cmd) {
	return m, nil
}

func AskMulti[T any](question string, choices []*Choice[T]) []*Choice[T] {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}
	list := MultiSelectModel[T]{
		items:    choices,
		selected: make(map[int]struct{}),
	}
	cursor := NewListCtrl[*MultiSelectModel[T]](&list)
	cursor.question = question
	cursor.debug = dump

	p := tea.NewProgram(cursor, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error while making selection: %v", err)
		os.Exit(1)
	}

	selections := make([]*Choice[T], 0, cursor.list.GetSelectedCount())
	for id, _ := range list.selected {
		selections = append(selections, choices[id])
	}

	return selections
}
