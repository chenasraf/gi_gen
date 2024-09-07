package ui

import (
	"strings"

	"github.com/rivo/uniseg"
)

const (
	CHECKMARK = "✔"
	CURSOR    = "▶"
)

func viewRow[T any](active bool, selected bool, checkbox bool, i int, choice *Choice[T]) (string, int) {
	var (
		s       strings.Builder
		tok     Colored
		c       int
		cursor  = " "
		checked = " "
	)

	if active {
		cursor = CURSOR
		// cursor_ = ">"
	}
	if selected {
		// checked = "x"
		// v
		checked = CHECKMARK
	}

	s.WriteString(cursor)
	c += uniseg.GraphemeClusterCount(cursor)

	if checkbox {
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

	if selected {
		tok = NewColored(ColorReset, choice.Label)
	} else {
		tok = NewColored(ColorDim, choice.Label)
	}
	s.WriteString(tok.String())
	c += tok.TextLength()

	return s.String(), c
}

type ViewWindow[T any] struct {
	maxWidth     int
	height       int
	offset       int
	isSelected   func(int) bool
	cursor       int
	lastMovement int
	choices      []*Choice[T]
	checkbox     bool
	border       bool
}
