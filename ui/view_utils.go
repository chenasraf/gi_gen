package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	// "unicode/utf8"

	"github.com/chenasraf/utils"
	"github.com/rivo/uniseg"
)

func viewRow[T any](active bool, selected bool, checkbox bool, i int, choice *Choice[T]) (string, int) {
	cursor_ := " "
	if active {
		cursor_ = "▶︎"
		// cursor_ = ">"
	}
	checked := " "
	if selected {
		checked = "x"
	}

	var s strings.Builder
	txt, c, cc := "", 0, 0

	s.WriteString(cursor_)
	c += uniseg.GraphemeClusterCount(cursor_)

	if checkbox {
		txt, cc = BlueText("[")
		s.WriteString(txt)
		c += cc

		s.WriteString(checked)
		c += uniseg.GraphemeClusterCount(checked)

		txt, cc = BlueText("] ")
		s.WriteString(txt)
		c += cc
	}

	if selected {
		txt, cc = choice.Label, uniseg.GraphemeClusterCount(choice.Label)
	} else {
		txt, cc = DimText(choice.Label)
	}
	s.WriteString(txt)
	c += cc
	// _ = cc
	// c -= 7

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

var termWidth int

func setTermWidth() {
	if termWidth > 0 {
		return
	}
	res, err := utils.RunCmd("tput", "cols")
	if err == nil {
		w, err := strconv.Atoi(strings.TrimSpace(res))
		if err == nil {
			termWidth = w
		} else {
			termWidth = 80
		}
	} else {
		termWidth = 80
	}
}

func viewRowsWindow[T any](w ViewWindow[T]) string {
	setTermWidth()
	var s strings.Builder
	isCursorAtStartEdge := w.cursor < w.offset
	isCursorAtEndEdge := w.cursor > len(w.choices)-w.height+w.offset
	startOffset := 0
	maxWidth := w.maxWidth
	if maxWidth == 0 {
		maxWidth = 80
	}
	width := int(math.Min(float64(termWidth), float64(w.maxWidth)))
	border := w.border
	if isCursorAtEndEdge {
		startOffset = len(w.choices) - w.cursor - w.height + w.offset
	}
	if isCursorAtStartEdge {
		startOffset = int(math.Abs(float64(w.cursor - w.offset)))
	}
	// s.WriteString(
	// 	fmt.Sprintf("lastMove: %d, offset: %d, startOffset: %d, cur: %d, len: %d, atEnd: %v\n", w.lastMovement, w.offset, startOffset, w.pos, len(w.choices), isCursorAtEndEdge),
	// )
	if border {
		s.WriteString("┌")
		s.WriteString(strings.Repeat("─", width-2))
		s.WriteString("┐\n")
	}
	for row := range w.height {
		i := w.cursor - w.offset + startOffset + row
		selected := w.isSelected(i)
		active := w.cursor == i
		choice := w.choices[i]
		rowTxt, contentLen := viewRow(active, selected, w.checkbox, i, choice)
		rowLen := contentLen
		if border {
			s.WriteString("│")
			// border + right padding
			rowLen += 3
		}
		s.WriteString(rowTxt)

		// dbg := fmt.Sprintf("w: %d, rl: %d, cl: %d", width, rowLen, contentLen)
		dbg := fmt.Sprintf("w: %d, rl: %d", width, rowLen)
		// dbg := fmt.Sprintf("")
		rowLen += uniseg.GraphemeClusterCount(dbg)

		if border {
			s.WriteString(strings.Repeat(" ", width-rowLen))
			s.WriteString(dbg)
			s.WriteString(" │")
		} else if len(dbg) > 0 {
			s.WriteString(strings.Repeat(" ", width-rowLen))
			s.WriteString(dbg)
		}

		s.WriteString("\n")
	}
	if border {
		s.WriteString("└")
		s.WriteString(strings.Repeat("─", width-2))
		s.WriteString("┘\n")
	}
	return s.String()
}
