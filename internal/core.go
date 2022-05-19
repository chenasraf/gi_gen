package internal

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func HandleFileOverwrite(path string, contents string) {
	overwriteSelection := AskOverwrite()
	switch overwriteSelection {
	case "":
		Quit()
		break
	case "Overwrite":
		log.Printf("Writing to %s", path)
		WriteFile(path, contents, true)
		break
	case "Append":
		log.Printf("Appending to %s", path)
		WriteFile(path, contents, false)
		break
	}
}

func AskOverwrite() string {
	overwritePrompt := &survey.Select{
		Message: ".gitignore file found in this directory. Please pick an option:",
		Options: []string{"Overwrite", "Append", "Skip"},
	}
	overwriteSelection := ""
	survey.AskOne(overwritePrompt, &overwriteSelection)
	return overwriteSelection
}

func GetCleanupSelection() bool {
	cleanupPrompt := &survey.Confirm{
		Message: "Do you want to remove patterns not existing in your project?",
		Default: true,
	}

	var cleanupSelection bool
	survey.AskOne(cleanupPrompt, &cleanupSelection)
	return cleanupSelection
}

func AskLanguage(fileNames []string, selected []string, files map[string]string) ([]string, []string) {
	langPrompt := &survey.MultiSelect{
		Message: "Found " + fmt.Sprint(len(fileNames)) +
			" possible matches in your project for gitignore files.\n" +
			"Please select which you want to write to .gitignore:\n",
		Options: fileNames,
	}

	var langSelections []string
	survey.AskOne(langPrompt, &langSelections)
	if langSelections == nil {
		Quit()
	}
	keys := []string{}
	for _, selection := range langSelections {
		selected = append(selected, files[selection])
		keys = append(keys, selection)
	}

	return selected, keys
}

func AskDiscovery() bool {
	prompt := &survey.Confirm{
		Message: "Would you like to try to scan for available templates automatically?\n" +
			"Select 'No' ('n') to see all available templates",
		Default: true,
	}
	var answer bool
	survey.AskOne(prompt, &answer)
	return answer
}

func Quit() {
	os.Exit(1)
}