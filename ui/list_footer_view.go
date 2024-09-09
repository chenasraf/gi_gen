package ui

import (
	"fmt"
	"math"
	"strings"

	utils "github.com/chenasraf/goutils"
)

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

	s.WriteString(wrapWithBorder(border, "", 0, width))

	for i, key := range m.ActiveHelpKeys() {
		str, cc := renderHelpKey(key)

		nextWillOverflow := (i+1)%fit == 0

		if !nextWillOverflow {
			spacing := width/fit - cc
			sep := strings.Repeat(" ", utils.MaxInt(0, spacing))
			str += sep
			line += str
			c += cc + len(sep)
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

func (m *ListCtrl[T, C]) FooterHeight() int {
	rowCount := m.GetHelpRowCount()

	h := 0
	if m.list.Style().HelpEnabled {
		h += rowCount + 1
	}

	return h
}
