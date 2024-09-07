package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type FilterMode int

const (
	FilterInactive = iota
	FilterFocused
	FilterActive
)

type ListFilter[T any] struct {
	term string
	mode FilterMode
}

func (m *ListFilter[T]) Update(msg tea.Msg) (*ListFilter[T], tea.Cmd) {
	spew.Fprintf(debug, "Filter Update: %v\n", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Batch(m.SetInactive, m.ClearTerm)
		case "enter":
			return m, m.SetActive
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			return m, m.TermhRemove(1)
		case "ctrl+w":
			if m.term == "" {
				return m, nil
			}
			split := strings.Split(m.term, " ")
			for split[len(split)-1] == "" {
				split = split[0 : len(split)-1]
			}
			lastWordLen := len(split[len(split)-1])
			if len(split) > 1 {
				lastWordLen++
			}
			return m, m.TermhRemove(lastWordLen)
		default:
			if len(msg.String()) == 1 {
				return m, m.TermAppend(msg.String())
			}
		}
	}
	return m, nil
}

func (m *ListFilter[T]) Filter(items []*Choice[T]) []*Choice[T] {
	var results []*Choice[T]
	for _, choice := range items {
		if m.IsVisible(choice) {
			results = append(results, choice)
		}
	}
	return results
}

func (m *ListFilter[T]) IsVisible(choice *Choice[T]) bool {
	t := strings.TrimSpace(m.term)
	t = strings.ToLower(t)
	l := strings.TrimSpace(choice.Label)
	l = strings.ToLower(l)
	if strings.Contains(l, t) {
		return true
	}
	return false
}

type FilterModeChange int

func (m *ListFilter[T]) SetFocused() tea.Msg {
	m.mode = FilterFocused
	return FilterModeChange(m.mode)
}

func (m *ListFilter[T]) SetInactive() tea.Msg {
	m.mode = FilterInactive
	return FilterModeChange(m.mode)
}

func (m *ListFilter[T]) SetActive() tea.Msg {
	m.mode = FilterActive
	return FilterModeChange(m.mode)
}

type FilterTermChange string

func (m *ListFilter[T]) ClearTerm() tea.Msg {
	m.term = ""
	return FilterTermChange(m.term)
}

func (m *ListFilter[T]) TermAppend(chr string) tea.Cmd {
	return func() tea.Msg {
		m.term += chr
		return FilterTermChange(m.term)
	}
}

func (m *ListFilter[T]) TermhRemove(n int) tea.Cmd {
	return func() tea.Msg {
		if len(m.term) > 0 {
			m.term = m.term[:len(m.term)-n]
		}
		return FilterTermChange(m.term)
	}
}
