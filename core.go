package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func handleFileOverwrite(path string, contents string) {
	overwriteSelection := getOverwriteSelection()
	switch overwriteSelection {
	case "":
		quit()
		break
	case "Overwrite":
		log.Printf("Writing to %s", path)
		writeFile(path, contents, true)
		break
	case "Append":
		log.Printf("Appending to %s", path)
		writeFile(path, contents, false)
		break
	}
}

func getOverwriteSelection() string {
	overwritePrompt := &survey.Select{
		Message: ".gitignore file found in this directory. Please pick an option:",
		Options: []string{"Overwrite", "Append", "Skip"},
	}
	overwriteSelection := ""
	survey.AskOne(overwritePrompt, &overwriteSelection)
	return overwriteSelection
}

func getCleanupSelection() bool {
	cleanupPrompt := &survey.Confirm{
		Message: "Do you want to remove patterns not existing in your project?",
		Default: true,
	}

	var cleanupSelection bool
	survey.AskOne(cleanupPrompt, &cleanupSelection)
	return cleanupSelection
}

func getLanguageSelections(fileNames []string, selected []string, files map[string]string) ([]string, []string) {
	langPrompt := &survey.MultiSelect{
		Message: "Found " + fmt.Sprint(len(fileNames)) +
			" possible matches in your project for gitignore files.\n" +
			"Please select which you want to write to .gitignore:\n",
		Options: fileNames,
	}

	var langSelections []string
	survey.AskOne(langPrompt, &langSelections)
	if langSelections == nil {
		quit()
	}
	keys := []string{}
	for _, selection := range langSelections {
		selected = append(selected, files[selection])
		keys = append(keys, selection)
	}

	return selected, keys
}

func getPossibleFiles(allFiles []string, files map[string]string, fileNames []string) []string {
	for _, filename := range allFiles {
		contents := readFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		if findFileMatches(contents) {

			files[langName] = contents
			fileNames = append(fileNames, langName)
		}
	}
	return fileNames
}

func quit() {
	os.Exit(1)
}
