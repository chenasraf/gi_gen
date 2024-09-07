package ui

import (
	"fmt"

	"github.com/rivo/uniseg"
)

type Colored struct {
	text  string
	color int
}

func NewColored(color int, text string) Colored {
	return Colored{color: color, text: text}
}

func (c Colored) TextLength() int {
	return uniseg.GraphemeClusterCount(c.text)
}

func (c Colored) TokenLength() int {
	if c.color == ColorReset {
		return 0
	}
	return len(ColorToken(c.color) + ColorToken(ColorReset))
}

func (c Colored) String() string {
	return ColorText(c.color, c.text)
}

func ColorToken(color int) string {
	c := fmt.Sprintf("\x1B[%dm", color)
	return c
}

func ColorText(color int, text string) string {
	if color == ColorReset {
		return text
	}
	colored := ColorToken(color)
	reset := ColorToken(ColorReset)

	out := fmt.Sprintf("%s%s%s", colored, text, reset)

	return out
}

const (
	ColorRed     = 31
	ColorGreen   = 32
	ColorYellow  = 33
	ColorBlue    = 34
	ColorMagenta = 35
	ColorCyan    = 36
	ColorWhite   = 37
	ColorBlack   = 30
	ColorDim     = 2
	ColorReset   = 0
)
