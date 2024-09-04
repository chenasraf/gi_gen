package prompt

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chenasraf/utils"
)

type Choice[T any] struct {
	Label string
	Value T
}

func AskMulti[T any](question string, choices []*Choice[T]) []*Choice[T] {
	cursor := CursorM[T]{
		question: question,
		choices:  choices,
		cursor:   0,
		selected: make(map[int]struct{}),
	}

	p := tea.NewProgram(cursor)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error while making selection: %v", err)
		os.Exit(1)
	}

	selections := make([]*Choice[T], 0, len(cursor.selected))
	for id, _ := range cursor.selected {
		selections = append(selections, choices[id])
	}

	return selections
}

func AskSingle[T any](question string, choices []*Choice[T]) *Choice[T] {
	cursor := Cursor[T]{
		question: question,
		choices:  choices,
		cursor:   0,
		selected: -1,
	}

	p := tea.NewProgram(cursor)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error while making selection: %v", err)
		os.Exit(1)
	}

	if cursor.selected == -1 {
		return nil
	}

	return cursor.choices[cursor.selected]
}

func AskLanguages(question string, languages []string) []string {
	choices := utils.MapSlice(languages, func(l string) *Choice[string] {
		return &Choice[string]{Label: l, Value: l}
	})
	res := AskMulti[string](question, choices)
	return utils.MapSlice(res, func(c *Choice[string]) string { return c.Value })
}
