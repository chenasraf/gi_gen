package ui

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenasraf/goutils"
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
	// TODO help delegates
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

func (m *ListCtrl[T, C]) RenderList() string {
	var s strings.Builder
	height := m.ListHeight()
	choices := m.list.Items()
	border := m.list.Style().DrawBorder
	checkbox := m.list.Style().DrawCheckbox

	// Edge offset - no scroll when closer than X to the start/end edge
	edgeOffset := m.list.Style().EdgeOffset
	if edgeOffset == -1 {
		edgeOffset = int(math.Floor(float64(m.height) / 2))
	}

	endOffset := 0 // NOTE fixes end offset w/ or w/o border
	offset := 0
	width := m.width
	hasQuestion := len(m.question) > 0

	// NOTE prevents list visual overflow
	if hasQuestion {
		endOffset++
	}
	if edgeOffset%2 != 0 {
		edgeOffset--
	}

	// Sticky behavior - lock scroll when close to list edges
	startEdge := edgeOffset - (m.HeaderHeight() - 1)
	endEdge := len(choices) - edgeOffset + m.FooterHeight()
	isCursorAtStartEdge := m.cursor < startEdge
	isCursorAtEndEdge := m.cursor >= endEdge+-1

	if isCursorAtEndEdge {
		offset = m.cursor - endEdge + 2
	}
	botOffset := 0
	if !isCursorAtStartEdge && edgeOffset%2 == 0 {
		botOffset = 1
	}

	spew.Fprintf(debug, "cur: %d, atStartEdge %s, atEndEdge %s, offset %d\n", m.cursor, isCursorAtStartEdge, isCursorAtEndEdge, offset)

	min := int(math.Max(0, float64(m.cursor)-float64(height)/2))
	max := int(math.Max(float64(m.cursor)+float64(height)/2, float64(height))) + 1

	// Row iteration
	for row := range max - min {
		i := min + row - offset + botOffset
		if i >= len(m.items) {
			s.WriteString(wrapWithBorder(border, "", 0, width))
			continue
		}
		choice := m.items[i]
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
	return s.String()
}

func (m *ListCtrl[T, C]) RenderHeader() string {
	hasQuestion := len(m.question) > 0
	border := m.list.Style().DrawBorder
	var s strings.Builder
	if hasQuestion {
		s.WriteString(wrapWithBorder(border, m.question, utils.StrLen(m.question), m.width))
	}

	status, statusLen := m.RenderStatusBar()
	s.WriteString(wrapWithBorder(border, status, statusLen, m.width))
	s.WriteString(wrapWithBorder(border, "", 0, m.width))
	return s.String()
}

func (m *ListCtrl[T, C]) RenderStatusBar() (string, int) {
	var (
		s   strings.Builder
		tok Colored
		c   int
	)
	txt := ""
	color := ColorDim

	if m.filter.mode == FilterFocused {
		color = ColorYellow
		txt += fmt.Sprintf("Type to filter: %s_", m.filter.term)
	} else {
		txt += fmt.Sprintf("%d selected | %d visible", m.list.GetSelectedCount(), len(m.items))
	}
	if m.filter.mode == FilterActive {
		txt += fmt.Sprintf(" | %d total | ", len(m.list.Items()))
	}

	tok = NewColored(color, txt)
	s.WriteString(tok.String())
	c += tok.TextLength()

	if m.filter.mode == FilterActive {
		txt = "Filtering: "
		tok = NewColored(ColorDim, txt)
		s.WriteString(tok.String())
		c += tok.TextLength()

		txt = m.filter.term
		tok = NewColored(ColorYellow, txt)
		s.WriteString(tok.String())
		c += tok.TextLength()
	}

	return s.String(), c
}

type HelpEntry struct {
	key, action string
}

func NewHelpEntry(key, action string) HelpEntry {
	return HelpEntry{key, action}
}

func renderHelpKey(entry HelpEntry) (string, int) {
	var (
		tok Colored
		c   int
	)
	tok = NewColored(ColorWhite, entry.key)
	keyStr := tok.String()
	c += tok.TextLength()

	tok = NewColored(ColorDim, entry.action)
	actionStr := tok.String()
	c += tok.TextLength()
	c += 2
	return fmt.Sprintf("%s  %s", keyStr, actionStr), c
}

func (m *ListCtrl[T, C]) GetHelpColumnCount() int {
	max := 0
	for _, key := range m.ActiveHelpKeys() {
		if _, cc := renderHelpKey(key); cc > max {
			max = cc
		}
	}
	width := m.width
	if m.list.Style().DrawBorder {
		width -= 2
	}
	// NOTE 2 spaces min between items
	max += 2

	fit := width / max
	spew.Fprintf(debug, "[ColumnCount] width: %d, max: %d, fit: %d\n", width, max, fit)
	return utils.MaxInt(fit, 1)
}

func (m *ListCtrl[T, C]) GetHelpRowCount() int {
	fit := m.GetHelpColumnCount()
	return int(math.Ceil(float64(len(m.listHelpKeys)) / float64(fit)))
}

func (m *ListCtrl[T, C]) RenderFooter() string {
	border := m.list.Style().DrawBorder
	width := m.width
	var s strings.Builder

	var (
		lines   []string
		lengths []int
	)
	c := 0
	line := ""
	fit := m.GetHelpColumnCount()
	spew.Fprintf(debug, "fit: %d, width: %d\n", fit, width)

	s.WriteString(wrapWithBorder(border, "", 0, width))

	for i, key := range m.ActiveHelpKeys() {
		str, cc := renderHelpKey(key)

		nextWillOverflow := (i+1)%fit == 0

		if !nextWillOverflow {
			spacing := width/fit - cc
			spew.Fprintf(debug, "width: %d, fit: %d, c: %d, spacing: %d\n", width, fit, cc, spacing)
			sep := strings.Repeat(" ", utils.MaxInt(0, spacing))
			str += sep
			line += str
			c += cc + len(sep)
			spew.Fprintf(debug, "adding line: %s, c: %d, len: %d\n", str, c, len(line))
		} else {
			line += str
			c += cc
			lines = append(lines, line)
			lengths = append(lengths, c)
			line = ""
			c = 0
		}
	}

	if len(line) > 0 {
		lines = append(lines, line)
		lengths = append(lengths, c)
	}
	for i, line := range lines {
		spew.Fprintf(debug, "printing line: %s, len: %d\n", line, lengths[i])
		s.WriteString(wrapWithBorder(border, line, lengths[i], width))
	}

	return s.String()
}

func (m *ListCtrl[T, C]) ActiveHelpKeys() []HelpEntry {
	if m.filter.mode == FilterFocused {
		return m.filterHelpKeys
	}
	return m.listHelpKeys
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

func (m *ListCtrl[T, C]) HeaderHeight() int {
	h := 3
	if len(m.question) == 0 {
		h--
	}
	if !m.list.Style().StatusEnabled {
		h--
	}
	return 3
}

func (m *ListCtrl[T, C]) FooterHeight() int {
	rowCount := m.GetHelpRowCount()

	h := 0
	if m.list.Style().HelpEnabled {
		h += rowCount + 1
	}

	return h
}

func (m *ListCtrl[T, C]) ListHeight() int {
	return m.height - m.HeaderHeight() - m.FooterHeight()
}
