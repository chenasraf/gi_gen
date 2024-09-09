package ui

import (
	"io"
	"math"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type ListModel[T any] interface {
	// TODO maybe return tea.Cmd in these so more behavior can be added
	Items() []*Choice[T]
	Select(choice *Choice[T])
	ClearSelected()
	SelectAll()
	GetSelectedCount() int
	IsSelected(choice *Choice[T]) bool
	Style() ListStyle
	Update(msg tea.Msg) (ListModel[T], tea.Cmd)
}

type ListStyle struct {
	// offset from/until which to scroll the list with
	// closer to the edge than that, and scrolling is stuck
	EdgeOffset    int
	DrawBorder    bool
	DrawCheckbox  bool
	FilterEnabled bool
	StatusEnabled bool
	HelpEnabled   bool
}

func DefaultListStyle() ListStyle {
	return ListStyle{
		EdgeOffset:    -1,
		FilterEnabled: true,
		StatusEnabled: true,
		HelpEnabled:   true,
	}
}

type ExitState int

const (
	ExitNormal = iota
	ExitQuit
)

var debug io.Writer

type ListCtrl[T ListModel[C], C comparable] struct {
	question       string
	width, height  int
	cursor         int
	list           T
	exitState      ExitState
	filter         *ListFilter[C]
	items          []*Choice[C]
	listHelpKeys   []HelpEntry
	filterHelpKeys []HelpEntry
}

const (
	BORDER_TOP_LEFT     = "┌"
	BORDER_BOTTOM_LEFT  = "└"
	BORDER_TOP_RIGHT    = "┐"
	BORDER_BOTTOM_RIGHT = "┘"
	BORDER_HORIZONTAL   = "─"
	BORDER_VERTICAL     = "│"
)

func NewListCtrl[T ListModel[C], C comparable](list T) *ListCtrl[T, C] {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}
	debug = dump
	filter := &ListFilter[C]{}
	listHelpKeys := []HelpEntry{
		NewHelpEntry("j/k/up/down", "move"),
		NewHelpEntry("space", "toggle"),
		NewHelpEntry("enter", "confirm"),
		NewHelpEntry("/", "search"),
		// NewHelpEntry("a", "select all"),
		// NewHelpEntry("c", "clear selection"),
	}
	filterHelpKeys := []HelpEntry{
		NewHelpEntry("esc", "clear filter"),
		NewHelpEntry("enter", "confirm"),
	}

	return &ListCtrl[T, C]{
		width:          80,
		height:         40,
		list:           list,
		filter:         filter,
		items:          list.Items(),
		listHelpKeys:   listHelpKeys,
		filterHelpKeys: filterHelpKeys,
	}
}

func (m *ListCtrl[T, C]) View() string {
	var s strings.Builder

	border := m.list.Style().DrawBorder

	// Edge offset - no scroll when closer than X to the start/end edge
	edgeOffset := m.list.Style().EdgeOffset
	if edgeOffset == -1 {
		edgeOffset = int(math.Floor(float64(m.height) / 2))
	}

	endOffset := 0 // NOTE fixes end offset w/ or w/o border
	width := m.width
	hasQuestion := len(m.question) > 0

	// NOTE prevents list visual overflow
	if hasQuestion {
		endOffset++
	}
	if edgeOffset%2 != 0 {
		edgeOffset--
	}

	// Top border
	if border && width > 1 {
		s.WriteString(BORDER_TOP_LEFT)
		s.WriteString(strings.Repeat(BORDER_HORIZONTAL, width-2))
		s.WriteString(BORDER_TOP_RIGHT + "\n")
	}

	s.WriteString(m.RenderHeader())
	s.WriteString(m.RenderList())
	s.WriteString(m.RenderFooter())

	// Bottom border
	if border && width > 1 {
		s.WriteString(BORDER_BOTTOM_LEFT)
		s.WriteString(strings.Repeat(BORDER_HORIZONTAL, width-2))
		s.WriteString(BORDER_BOTTOM_RIGHT)
	}

	return s.String()
}

func (m *ListCtrl[T, C]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	spew.Fprintf(debug, "msg: %v\n", msg)

	if m.filter.mode == FilterFocused {
		switch msg := msg.(type) {
		case FilterTermChange:
			spew.Fprintf(debug, "FilterTermChange: %v\n", msg)
			items := m.list.Items()
			m.items = m.filter.Filter(items)
			if m.cursor >= len(m.items) {
				cmds = append(cmds, m.MoveCursor(0))
			}
		}
		filter, cmd := m.filter.Update(msg)
		m.filter = filter
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
		rem := 1
		if m.list.Style().DrawBorder {
			rem += 2
		}
		// TODO include header, footer
		m.SetHeight(msg.Height - rem)
	case FilterModeChange:
		spew.Fprintf(debug, "FilterModeChange: %v\n", msg)
		if msg == FilterInactive {
			m.items = m.list.Items()
		}
		cmds = append(cmds, m.MoveCursor(0))
	case tea.KeyMsg:

		switch msg.String() {
		// TODO use bound config keys
		case "ctrl+c", "q":
			m.exitState = 1
			return m, tea.Quit
		case "up", "k":
			return m, m.MoveCursor(-1)
		case "down", "j":
			return m, m.MoveCursor(1)
		case "left", "pgup", "h":
			return m, m.MoveCursor(-10)
		case "right", "pgdown", "l":
			return m, m.MoveCursor(10)
		case " ":
			choice := m.items[m.cursor]
			return m, m.Select(choice)
		case "esc":
			if m.filter.mode == FilterActive {
				return m, tea.Batch(m.filter.SetInactive, m.filter.ClearTerm)
			}
		case "a":
			return m, m.SelectAll
		case "c":
			return m, m.ClearSelected
		case "enter":
			m.exitState = 0
			return m, tea.Quit
		case "/":
			return m, m.filter.SetFocused
		}
	}

	newList, cmd := m.list.Update(msg)
	m.list = newList.(T)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ListCtrl[T, C]) Init() tea.Cmd {
	return nil
}

type ListUpdateMsg int

func (m *ListCtrl[T, C]) Select(choice *Choice[C]) tea.Cmd {
	return func() tea.Msg {
		m.list.Select(choice)
		return ListUpdateMsg(m.list.GetSelectedCount())
	}
}

func (m *ListCtrl[T, C]) ClearSelected() tea.Msg {
	m.list.ClearSelected()
	return ListUpdateMsg(0)
}

func (m *ListCtrl[T, C]) SelectAll() tea.Msg {
	m.list.SelectAll()
	return ListUpdateMsg(len(m.list.Items()))
}

func (m *ListCtrl[T, C]) MoveCursor(amount int) tea.Cmd {
	m.cursor += amount
	size := len(m.items)
	if m.cursor >= size {
		m.cursor = m.cursor - size
	}
	if m.cursor < 0 {
		m.cursor = m.cursor + size
	}
	// second check, list size changed - force within bounds
	if m.cursor >= size {
		m.cursor = size - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	return nil
}

func (m *ListCtrl[T, C]) SetWidth(width int) {
	m.width = width
}

func (m *ListCtrl[T, C]) SetHeight(height int) {
	m.height = height
}

func wrapWithBorder(border bool, s string, size int, width int) string {
	if !border || width < 2 {
		return s + "\n"
	}

	var out strings.Builder
	out.WriteString(BORDER_VERTICAL + " ")
	out.WriteString(s)
	out.WriteString(strings.Repeat(" ", width-size-4)) // NOTE 4 = border+padding
	out.WriteString(" " + BORDER_VERTICAL + "\n")

	return out.String()
}
