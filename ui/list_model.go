package ui

import (
	"io"
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

type ListModel[T any] interface {
	// TODO maybe return tea.Cmd in these so more behavior can be added
	Items() []*Choice[T]
	Select(i int)
	GetSelectedCount() int
	IsSelected(i int) bool
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

type ListCtrl[T ListModel[C], C any] struct {
	question      string
	status        string
	width, height int
	cursor        int
	list          T
	exitState     ExitState
	debug         io.Writer
	filter        *ListFilter
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
	filter := &ListFilter{}
	return &ListCtrl[T, C]{list: list, filter: filter}
}

func (m *ListCtrl[T, C]) SetWidth(width int) {
	m.width = width
}

func (m *ListCtrl[T, C]) SetHeight(height int) {
	m.height = height
}

func (m *ListCtrl[T, C]) View() string {
	var s strings.Builder

	choices := m.list.Items()
	border := m.list.Style().DrawBorder
	checkbox := m.list.Style().DrawCheckbox

	// Edge offset - no scroll when closer than X to the start/end edge
	edgeOffset := m.list.Style().EdgeOffset
	if edgeOffset == -1 {
		edgeOffset = int(m.height / 2)
	}

	// Fix var for end offset w/ or w/o border
	endOffset := 0
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
		s.WriteString(wrapWithBorder(border, m.question, len(m.question), width))
		s.WriteString(wrapWithBorder(border, "", 0, width))
		height -= 1
	}

	spew.Fprintf(m.debug,
		"cursor: %d, offset: %d, startOffset: %d, len: %d, height: %d, atEnd: %v\n",
		m.cursor, edgeOffset, offset, len(choices), m.height, isCursorAtEndEdge,
	)

	// Row iteration
	for row := range height {
		i := m.cursor - edgeOffset + offset + row
		selected := m.list.IsSelected(i)
		active := m.cursor == i
		choice := choices[i]

		// TODO create & use row delegate
		rowTxt, contentLen := viewRow(active, selected, checkbox, i, choice)
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
	spew.Fprintf(m.debug, "msg: %v\n", msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetWidth(msg.Width)
		rem := 1
		if m.list.Style().DrawBorder {
			rem += 2
		}
		// TODO include header, footer
		m.SetHeight(msg.Height - rem)
	case tea.KeyMsg:

		if m.filter.mode == FilterFocused {
			filter, cmd := m.filter.Update(msg)
			m.filter = filter
			return m, cmd
		}

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
			return m, m.Select(m.cursor)
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

func (m *ListCtrl[T, C]) Select(i int) tea.Cmd {
	m.list.Select(i)
	return nil
}

func (m *ListCtrl[T, C]) MoveCursor(amount int) tea.Cmd {
	m.cursor += amount
	if m.cursor >= len(m.list.Items()) {
		m.cursor = m.cursor - len(m.list.Items())
	}
	if m.cursor < 0 {
		m.cursor = m.cursor + len(m.list.Items())
	}
	return nil
}
