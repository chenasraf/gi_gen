package ui

import tea "github.com/charmbracelet/bubbletea"

type FilterMode int

const (
	FilterInactive = iota
	FilterFocused
	FilterActive
)

type ListFilter struct {
	term string
	mode FilterMode
}

func (m *ListFilter) Update(msg tea.Msg) (*ListFilter, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Batch(m.SetInactive, m.Clear)
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m, m.SearchAppend(msg.String())
		}
	}
	return m, nil
}

func (m *ListFilter) SetFocused() tea.Msg {
	m.mode = FilterFocused
	return nil
}

func (m *ListFilter) SetInactive() tea.Msg {
	m.mode = FilterInactive
	return nil
}

func (m *ListFilter) Clear() tea.Msg {
	m.term = ""
	return nil
}

func (m *ListFilter) SearchAppend(chr string) tea.Cmd {
	return func() tea.Msg {
		m.term += chr
		return nil
	}
}
