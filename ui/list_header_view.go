package ui

import (
	"fmt"
	"strings"

	utils "github.com/chenasraf/goutils"
)

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
