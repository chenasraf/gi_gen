package main

import (
	"errors"
	"fmt"

	"github.com/chenasraf/gi_gen/prompt"
	"github.com/chenasraf/utils"
)

func main() {
	opts := parseArgs()

	if opts.autoSelect {
		// TODO extract to languages.go
		languages, _ := getAutoSelectCandidates()
		if len(languages) == 1 {
			// Single candidate, auto-select
			*opts.languages = languages
		} else {
			question := ""
			var choices []string

			if len(languages) > 1 {
				// Multiple candidates, prompt for selection
				question = "Found more than one candidate! Please select which to use:"
				choices = languages
			} else {
				// No candidates, prompt for selection
				question = "Couldn't auto-detect project type. Please select gitignore templates:"
				allTemplates := getAllIgnoreTemplates()
				choices = utils.MapKeys(allTemplates)
			}

			*opts.languages = prompt.AskLanguages(question, choices)
		}
	} else {
		// TODO extract to languages.go
		choices := utils.MapKeys(getAllIgnoreTemplates())
		*opts.languages = prompt.AskLanguages("Please select gitignore templates to use:", choices)
	}

	fmt.Println("Selected languages:", *opts.languages)

	if len(*opts.languages) == 0 {
		utils.HandleErr(errors.New("No languages selected"))
	}

	if !opts.outputBehavior.IsValid() {
		// TODO handle default output behavior
	}

	if !opts.conflictBehavior.IsValid() {
		// TODO handle default conflict behavior
	}

	// TODO continue with processing
}
