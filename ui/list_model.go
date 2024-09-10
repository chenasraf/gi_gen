package ui

import (
	"io"
	"os"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	utils "github.com/chenasraf/goutils"
	"github.com/davecgh/go-spew/spew"
)

// ListModel is an interface that defines the methods required for a list model.
// It includes methods for managing items, selection, and updating the list.
type ListModel[T any] interface {
	// Items returns a slice of choices in the list.
	Items() []*Choice[T]
	// Select marks a choice as selected.
	Select(choice *Choice[T])
	// ClearSelected clears all selected choices.
	ClearSelected()
	// SelectAll selects all choices in the list.
	SelectAll()
	// GetSelectedCount returns the count of selected choices.
	GetSelectedCount() int
	// IsSelected checks if a choice is selected.
	IsSelected(choice *Choice[T]) bool
	// Style returns the style settings for the list.
	Style() ListStyle
	// Update handles messages and updates the list state.
	Update(msg tea.Msg) (ListModel[T], tea.Cmd)
}

// ListStyle defines the visual style and behavior settings for the list.
type ListStyle struct {
	// DrawBorder indicates whether to draw a border around the list.
	DrawBorder bool
	// DrawCheckbox indicates whether to draw checkboxes for selections.
	DrawCheckbox bool
	// FilterEnabled indicates whether filtering is enabled.
	FilterEnabled bool
	// StatusEnabled indicates whether status display is enabled.
	StatusEnabled bool
	// HelpEnabled indicates whether help display is enabled.
	HelpEnabled bool
}

// DefaultListStyle returns a ListStyle with default settings.
func DefaultListStyle() ListStyle {
	return ListStyle{
		FilterEnabled: true,
		StatusEnabled: true,
		HelpEnabled:   true,
	}
}

// ExitState represents the state of exiting the list control.
type ExitState int

const (
	// ExitNormal indicates a normal exit state.
	ExitNormal = iota
	// ExitQuit indicates an exit state due to quitting.
	ExitQuit
)

var debug io.Writer

// ListCtrl is a controller for managing a list of items with various functionalities
// such as filtering, selecting, and navigating through the list.
type ListCtrl[T ListModel[C], C comparable] struct {
	question       string
	width, height  int
	cursor         int
	list           T
	exitState      ExitState
	helpActive     bool
	filter         *ListFilter[C]
	items          []*Choice[C]
	listHelpKeys   []HelpEntry
	filterHelpKeys []HelpEntry
	fullHelpKeys   []HelpEntry
}

// Border characters used for drawing the list borders.
const (
	BORDER_TOP_LEFT     = "┌"
	BORDER_BOTTOM_LEFT  = "└"
	BORDER_TOP_RIGHT    = "┐"
	BORDER_BOTTOM_RIGHT = "┘"
	BORDER_HORIZONTAL   = "─"
	BORDER_VERTICAL     = "│"
)

// NewListCtrl creates a new ListCtrl instance with the given list model.
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
		NewHelpEntry("?", "help"),
	}
	filterHelpKeys := []HelpEntry{
		NewHelpEntry("esc", "clear filter"),
		NewHelpEntry("enter", "confirm"),
	}
	fullHelpKeys := []HelpEntry{
		NewHelpEntry("j/k/up/down", "move up/down"),
		NewHelpEntry("h/l/left/right/pgup/pgdown", "move page up/down"),
		NewHelpEntry("space", "toggle"),
		NewHelpEntry("enter", "confirm"),
		NewHelpEntry("/", "search"),
		NewHelpEntry("?", "close help"),
		NewHelpEntry("a", "select all"),
		NewHelpEntry("c", "clear selection"),
		NewHelpEntry("esc", "clear filter"),
	}

	return &ListCtrl[T, C]{
		width:          80,
		height:         40,
		list:           list,
		filter:         filter,
		items:          list.Items(),
		listHelpKeys:   listHelpKeys,
		filterHelpKeys: filterHelpKeys,
		fullHelpKeys:   fullHelpKeys,
	}
}

// View renders the list view as a string.
func (m *ListCtrl[T, C]) View() string {
	var s strings.Builder
	border := m.list.Style().DrawBorder
	width := m.width

	// Top border
	if border && width > 1 {
		s.WriteString(BORDER_TOP_LEFT)
		s.WriteString(strings.Repeat(BORDER_HORIZONTAL, width-2))
		s.WriteString(BORDER_TOP_RIGHT + "\n")
	}

	// Content
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

// Update handles messages and updates the list state accordingly.
func (m *ListCtrl[T, C]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	spew.Fprintf(debug, "msg: (%s) %v\n", reflect.TypeOf(msg), msg)

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
		case "?":
			return m, m.ToggleHelp
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

// ListUpdateMsg is a type used for messages that update the list.
type ListUpdateMsg int

type FooterUpdateMsg int

// Select selects a choice from the list and returns a command.
func (m *ListCtrl[T, C]) Select(choice *Choice[C]) tea.Cmd {
	return func() tea.Msg {
		m.list.Select(choice)
		return ListUpdateMsg(m.list.GetSelectedCount())
	}
}

// ClearSelected clears all selected choices from the list and returns a message.
func (m *ListCtrl[T, C]) ClearSelected() tea.Msg {
	m.list.ClearSelected()
	return ListUpdateMsg(0)
}

// SelectAll selects all choices in the list and returns a message.
func (m *ListCtrl[T, C]) SelectAll() tea.Msg {
	m.list.SelectAll()
	return ListUpdateMsg(len(m.list.Items()))
}

func (m *ListCtrl[T, C]) ToggleHelp() tea.Msg {
	m.helpActive = !m.helpActive
	return FooterUpdateMsg(utils.Ternary(m.helpActive, 1, 0))
}

// MoveCursor moves the cursor by the specified amount and returns a command.
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

// SetWidth sets the width of the list.
func (m *ListCtrl[T, C]) SetWidth(width int) {
	m.width = width
}

// SetHeight sets the height of the list.
func (m *ListCtrl[T, C]) SetHeight(height int) {
	m.height = height
}

// wrapWithBorder wraps a string with a border if the border is enabled.
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
