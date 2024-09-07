package ui

import (
	"github.com/chenasraf/utils"
)

type Choice[T any] struct {
	Label string
	Value T
}

func AskLanguages(question string, languages []string) []string {
	choices := utils.MapSlice(languages, func(l string) *Choice[string] {
		return &Choice[string]{Label: l, Value: l}
	})
	res := AskMulti[string](question, choices)
	return utils.MapSlice(res, func(c *Choice[string]) string { return c.Value })
}
