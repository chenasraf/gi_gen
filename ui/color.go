package ui

import "fmt"

func Color(color int) (string, int) {
	c := fmt.Sprintf("\x1B[%dm", color)
	return c, 0 // len(c)
}

func ColorText(color int, text string) (string, int) {
	colored, colorCount := Color(color)
	reset, resetCount := Reset()

	out := fmt.Sprintf("%s%s%s", colored, text, reset)
	count := (colorCount + resetCount + len(text))

	return out, count
}

func Red() (string, int) {
	return Color(31)
}

func RedText(text string) (string, int) {
	return ColorText(31, text)
}

func Green() (string, int) {
	return Color(32)
}

func GreenText(text string) (string, int) {
	return ColorText(32, text)
}

func Yellow() (string, int) {
	return Color(33)
}

func YellowText(text string) (string, int) {
	return ColorText(33, text)
}

func Blue() (string, int) {
	return Color(34)
}

func BlueText(text string) (string, int) {
	return ColorText(34, text)
}

func Magenta() (string, int) {
	return Color(35)
}

func MagentaText(text string) (string, int) {
	return ColorText(35, text)
}

func Cyan() (string, int) {
	return Color(36)
}

func CyanText(text string) (string, int) {
	return ColorText(36, text)
}

func White() (string, int) {
	return Color(37)
}

func WhiteText(text string) (string, int) {
	return ColorText(37, text)
}

func Black() (string, int) {
	return Color(30)
}

func BlackText(text string) (string, int) {
	return ColorText(30, text)
}

func Dim() (string, int) {
	return Color(2)
}

func DimText(text string) (string, int) {
	return ColorText(2, text)
}

func Reset() (string, int) {
	return Color(0)
}
