package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chenasraf/utils"
)

func main() {
	opts := parseArgs()
	isAuto := false

	*opts.languages, isAuto = getLanguageSelections(opts)

	if isAuto {
		fmt.Println("Detected languages:", strings.Join(*opts.languages, ", "))
	} else {
		fmt.Println("Selected languages:", strings.Join(*opts.languages, ", "))
	}

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
