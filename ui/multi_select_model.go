package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type MultiSelectModel[T comparable] struct {
	items    []*Choice[T]
	selected map[T]struct{}
}

func (m *MultiSelectModel[T]) Items() []*Choice[T] {
	return m.items
}

func (m *MultiSelectModel[T]) Select(choice *Choice[T]) {
	if _, exists := m.selected[choice.Value]; exists {
		delete(m.selected, choice.Value)
	} else {
		m.selected[choice.Value] = struct{}{}
	}
}

func (m *MultiSelectModel[T]) GetSelectedCount() int {
	return len(m.selected)
}

func (m *MultiSelectModel[T]) IsSelected(choice *Choice[T]) bool {
	_, exists := m.selected[choice.Value]
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

func AskMulti[T comparable](question string, choices []*Choice[T]) []*Choice[T] {
	list := MultiSelectModel[T]{
		items:    choices,
		selected: make(map[T]struct{}),
	}
	cursor := NewListCtrl[*MultiSelectModel[T]](&list)
	cursor.question = question

	p := tea.NewProgram(cursor, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error while making selection: %v", err)
		os.Exit(1)
	}

	selections := make([]*Choice[T], 0, cursor.list.GetSelectedCount())
	choiceMap := make(map[T]*Choice[T])
	for _, choice := range list.Items() {
		choiceMap[choice.Value] = choice
	}
	for value, _ := range list.selected {
		selections = append(selections, choiceMap[value])
	}

	return selections
}
