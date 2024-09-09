package ui

import (
	"strings"

	"github.com/chenasraf/goutils"
)

const (
	CHECKMARK = "✔"
	CURSOR    = "▶"
)

type RowModel[T any] struct {
	selected, checkbox, active bool
	choice                     *Choice[T]
}

func (m *RowModel[T]) Render() (string, int) {
	var (
		s       strings.Builder
		tok     Colored
		c       int
		cursor  = " "
		checked = " "
	)

	if m.active {
		cursor = CURSOR
	}
	if m.selected {
		checked = CHECKMARK
	}

	s.WriteString(cursor)
	c += utils.StrLen(cursor)

	if m.checkbox {
		tok = NewColored(ColorBlue, "[")
		s.WriteString(tok.String())
		c += tok.TextLength()

		tok = NewColored(ColorReset, checked)
		s.WriteString(tok.String())
		c += tok.TextLength()

		tok = NewColored(ColorBlue, "] ")
		s.WriteString(tok.String())
		c += tok.TextLength()
	}

	if m.selected {
		tok = NewColored(ColorReset, m.choice.Label)
	} else {
		tok = NewColored(ColorDim, m.choice.Label)
	}
	s.WriteString(tok.String())
	c += tok.TextLength()

	return s.String(), c
}
