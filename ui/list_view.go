package ui

import (
	"strings"

	utils "github.com/chenasraf/goutils"
	"github.com/davecgh/go-spew/spew"
)

func (m *ListCtrl[T, C]) RenderList() string {
	var s strings.Builder
	height := m.ListHeight()
	border := m.list.Style().DrawBorder
	checkbox := m.list.Style().DrawCheckbox
	width := m.width

	midF := float64(height) / 2
	mid := int(midF)
	min := 0
	max := len(m.items)

	offset := -1
	isNearEnd := m.cursor >= max-mid+offset

	if isNearEnd {
		min = len(m.items) - height + offset
	} else {
		min = int(utils.MaxInt(0, m.cursor-mid))
	}
	min = utils.MaxInt(0, min)
	max = min + height - offset

	// Row iteration
	for l := range max - min {
		row := min + l
		spew.Fprintf(debug, "row: %d, min: %d, max: %d, len: %d\n", row, min, max, len(m.items))
		if row >= len(m.items) {
			s.WriteString(wrapWithBorder(border, "", 0, width))
			continue
		}
		choice := m.items[row]
		active := m.cursor == row
		selected := m.list.IsSelected(choice)

		model := RowModel[C]{
			choice:   choice,
			selected: selected,
			active:   active,
			checkbox: checkbox,
		}
		rowTxt, contentLen := model.Render()
		s.WriteString(wrapWithBorder(border, rowTxt, contentLen, width))
	}
	return s.String()
}

func (m *ListCtrl[T, C]) ListHeight() int {
	return m.height - m.HeaderHeight() - m.FooterHeight()
}
