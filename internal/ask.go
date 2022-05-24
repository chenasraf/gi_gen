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
	return askYesNo("Would you like to try to scan for available templates automatically?\n" +
		"Select 'No' to see all available templates")
}

func AskOverwrite() string {
	fmt.Println()
	return askSelection(
		".gitignore file found in this directory. Please pick an option:",
		[]string{"Overwrite", "Append", "Skip"},
	)
}

func AskCleanup() bool {
	fmt.Println()
	return askYesNo("Do you want to remove patterns not existing in your project?\nThis might produce incomplete files on new projects.")
}

func askYesNo(message string) bool {
	prompt := &survey.Select{
		Message: message,
		Default: "Yes",
		Options: []string{"Yes", "No"},
	}
	answer := ""
	survey.AskOne(prompt, &answer)
	if answer == "" {
		KeyInterrupt()
	}

	return answer == "Yes"
}

func askMulti(message string, options []string) []string {
	langPrompt := &survey.MultiSelect{
		Message: message,
		Options: options,
	}

	var selections []string
	survey.AskOne(langPrompt, &selections)

	if selections == nil {
		KeyInterrupt()
	}

	return selections
}

func askSelection(message string, options []string) string {
	langPrompt := &survey.Select{
		Message: message,
		Options: options,
	}

	selection := ""
	survey.AskOne(langPrompt, &selection)

	if selection == "" {
		KeyInterrupt()
	}

	return selection
}
