package ui

import (
	"math"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

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

func (m *ListCtrl[T, C]) ListHeight() int {
	return m.height - m.HeaderHeight() - m.FooterHeight()
}
