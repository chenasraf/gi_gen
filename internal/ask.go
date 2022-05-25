package internal

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

func askLanguage(fileNames []string, selected []string, files map[string]string) ([]string, []string) {
	fmt.Println()
	langSelections := askMulti("Please select which you want to write to .gitignore:\n", maps.Keys(files))

	keys := []string{}
	for _, selection := range langSelections {
		selected = append(selected, files[selection])
		keys = append(keys, selection)
	}

	return selected, keys
}

func askDiscovery() bool {
	return askYesNo("Would you like to try to scan for available templates automatically?\n"+
		"Select 'No' to see all available templates", true)
}

func askOverwrite() string {
	fmt.Println()
	return askSelection(
		".gitignore file found in this directory. Please pick an option:",
		[]string{"Overwrite", "Append", "Skip"},
		"Overwrite",
	)
}

func askCleanup() bool {
	fmt.Println()
	return askYesNo("Do you want to remove patterns not existing in your project?\n"+
		"This might produce incomplete files on new projects.", false)
}

func askYesNo(message string, defaultValue bool) bool {
	return askSelection(message, []string{"Yes", "No"}, ternary(defaultValue, "Yes", "No")) == "Yes"
}

func askMulti(message string, options []string) []string {
	langPrompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}

	var selections []string
	survey.AskOne(langPrompt, &selections)

	if selections == nil {
		keyInterrupt()
	}

	return selections
}

func askSelection(message string, options []string, defaultValue string) string {
	langPrompt := &survey.Select{
		Message: message,
		Options: options,
		Default: defaultValue,
	}

	selection := ""
	survey.AskOne(langPrompt, &selection)

	if selection == "" {
		keyInterrupt()
	}

	return selection
}
