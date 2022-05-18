package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/exp/maps"
)

// create line cache map for skipping previously checked glob lines
var lineCache = make(map[string]bool)

func main() {
	allFiles, err := prepareGitignores()

	handleErr(err)

	files := make(map[string]string)
	fileNames := []string{}

	for _, filename := range allFiles {
		contents := readFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		if findFileMatches(contents) {
			sep := "#========================================================================\n"
			header := fmt.Sprintf(sep+"# %s\n"+sep+"\n", basename)
			files[langName] = header + contents
			fileNames = append(fileNames, langName)
		}
	}

	log.Println("Done.")

	selected := []string{}
	keys := maps.Keys(files)

	if len(keys) > 1 {
		prompt := &survey.MultiSelect{
			Message: "Found matches in your project for the following gitignore files.\nPlease select which you want to write to .gitignore:",
			Options: fileNames,
		}

		var inp []string
		survey.AskOne(prompt, &inp)

		for _, sel := range inp {
			selected = append(selected, files[sel])
		}
	} else {
		selected = keys
	}

	wd, _ := os.Getwd()
	out := filepath.Join(wd, ".gitignore")

	if fileExists(out) {
		prompt2 := &survey.Select{
			Message: ".gitignore file found in this directory. Please pick an option:",
			Options: []string{"Overwrite", "Append", "Skip"},
		}
		var inp2 string
		survey.AskOne(prompt2, &inp2)
		switch inp2 {
		case "Overwrite":
			log.Printf("Writing to %s", out)
			writeFile(out, strings.Join(selected, "\n"), true)
			break
		case "Append":
			log.Printf("Appending to %s", out)
			writeFile(out, strings.Join(selected, "\n"), false)
			break
		}
	} else {
		log.Printf("Writing to %s", out)
		writeFile(out, strings.Join(selected, "\n"), true)
	}
}
