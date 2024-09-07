package ui

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenasraf/utils"
	"github.com/davecgh/go-spew/spew"
)

type ListModel[T any] interface {
	// TODO maybe return tea.Cmd in these so more behavior can be added
	Items() []*Choice[T]
	Select(choice *Choice[T])
	GetSelectedCount() int
	IsSelected(choice *Choice[T]) bool
	Style() ListStyle
	Update(msg tea.Msg) (ListModel[T], tea.Cmd)
}

type ListStyle struct {
	// offset from/until which to scroll the list with
	// closer to the edge than that, and scrolling is stuck
	EdgeOffset   int
	DrawBorder   bool
	DrawCheckbox bool
}

func DefaultListStyle() ListStyle {
	return ListStyle{EdgeOffset: -1}
}

type ExitState int

const (
	ExitNormal = iota
	ExitQuit
)

var debug io.Writer

type ListCtrl[T ListModel[C], C any] struct {
	question      string
	status        string
	width, height int
	cursor        int
	list          T
	exitState     ExitState
	filter        *ListFilter[C]
	filteredItems []*Choice[C]
	// TODO help delegates
}

const (
	BORDER_LEFT         = "│"
	BORDER_RIGHT        = "│"
	BORDER_TOP_LEFT     = "┌"
	BORDER_TOP_RIGHT    = "┐"
	BORDER_BOTTOM_LEFT  = "└"
	BORDER_BOTTOM_RIGHT = "┘"
	BORDER_HORIZONTAL   = "─"
	BORDER_VERTICAL     = "│"
)

func NewListCtrl[T ListModel[C], C any](list T) *ListCtrl[T, C] {
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
	return &ListCtrl[T, C]{list: list, filter: filter, filteredItems: list.Items()}
}

func (m *ListCtrl[T, C]) SetWidth(width int) {
	m.width = width
}

func (m *ListCtrl[T, C]) SetHeight(height int) {
	m.height = height
}

func (m *ListCtrl[T, C]) View() string {
	var s strings.Builder

	spew.Fprintf(debug, "ListCtrl View: %s\n", m.filter.term)

	choices := m.list.Items()
	border := m.list.Style().DrawBorder
	checkbox := m.list.Style().DrawCheckbox

	// Edge offset - no scroll when closer than X to the start/end edge
	edgeOffset := m.list.Style().EdgeOffset
	if edgeOffset == -1 {
		edgeOffset = int(m.height / 2)
	}

	endOffset := 0 // NOTE fixes end offset w/ or w/o border
	offset := 0
	width := m.width
	height := m.height
	hasQuestion := len(m.question) > 0

	// NOTE prevents list visual overflow
	if hasQuestion {
		endOffset += 1
	}

	// Sticky behavior - lock scroll when close to list edges
	isCursorAtStartEdge := m.cursor < edgeOffset
	isCursorAtEndEdge := m.cursor > len(choices)-height+edgeOffset
	if isCursorAtStartEdge {
		offset = int(math.Abs(float64(m.cursor - edgeOffset)))
	}
	if isCursorAtEndEdge {
		offset = len(choices) - m.cursor - height + edgeOffset + endOffset
	}

	// Top border
	if border && width > 1 {
		s.WriteString(BORDER_TOP_LEFT)
		s.WriteString(strings.Repeat(BORDER_HORIZONTAL, width-2))
		s.WriteString(BORDER_TOP_RIGHT + "\n")
	}

	// Question + spacing
	if hasQuestion {
		s.WriteString(wrapWithBorder(border, m.question, utils.StrLen(m.question), width))
		height -= 1
	}

	status := m.status
	if m.filter.mode == FilterFocused {
		status = fmt.Sprintf("Type to filter: %s_", m.filter.term)
	}
	s.WriteString(wrapWithBorder(border, status, utils.StrLen(status), width))

	// spew.Fprintf(debug,
	// 	"cursor: %d, offset: %d, startOffset: %d, len: %d, height: %d, atEnd: %v\n",
	// 	m.cursor, edgeOffset, offset, len(choices), m.height, isCursorAtEndEdge,
	// )

	// Row iteration
	for row := range height {
		i := m.cursor - edgeOffset + offset + row
		if i >= len(m.filteredItems) {
			s.WriteString(wrapWithBorder(border, "", 0, width))
			continue
		}
		choice := m.filteredItems[i]
		active := m.cursor == i
		selected := m.list.IsSelected(choice)

		row := RowModel[C]{
			choice:   choice,
			selected: selected,
			active:   active,
			checkbox: checkbox,
		}
		rowTxt, contentLen := row.Render()
		rowLen := contentLen

		if border {
			rowLen += 3 // NOTE left/right borders + right padding
		}
		s.WriteString(wrapWithBorder(border, rowTxt, contentLen, width))
	}

	// Bottom border
	if border && width > 1 {
		s.WriteString(BORDER_BOTTOM_LEFT)
		s.WriteString(strings.Repeat(BORDER_HORIZONTAL, width-2))
		s.WriteString(BORDER_BOTTOM_RIGHT)
	}

	return s.String()
}

func wrapWithBorder(border bool, s string, size int, width int) string {
	if !border || width < 2 {
		return s + "\n"
	}

	var out strings.Builder
	out.WriteString(BORDER_LEFT + " ")
	out.WriteString(s)
	out.WriteString(strings.Repeat(" ", width-size-4)) // NOTE 4 = border+padding
	out.WriteString(" " + BORDER_RIGHT + "\n")

	return out.String()
}

func (m *ListCtrl[T, C]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	spew.Fprintf(debug, "msg: %v\n", msg)

	if m.filter.mode == FilterFocused {
		switch msg := msg.(type) {
		case FilterTermChange:
			spew.Fprintf(debug, "FilterTermChange: %v\n", msg)
			items := m.list.Items()
			m.filteredItems = m.filter.Filter(items)
			if m.cursor >= len(m.filteredItems) {
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
			m.filteredItems = m.list.Items()
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
		case "left", "pgup":
			return m, m.MoveCursor(-10)
		case "right", "pgdown":
			return m, m.MoveCursor(10)
		case " ":
			choice := m.filteredItems[m.cursor]
			return m, m.Select(choice)
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

func (m *ListCtrl[T, C]) Select(choice *Choice[C]) tea.Cmd {
	m.list.Select(choice)
	return nil
}

func (m *ListCtrl[T, C]) MoveCursor(amount int) tea.Cmd {
	m.cursor += amount
	size := len(m.filteredItems)
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
