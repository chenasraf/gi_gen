package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

func main() {
	wd, _ := os.Getwd()
	outFile := filepath.Join(wd, ".gitignore")
	allFiles, err := prepareGitignores()

	handleErr(err)

	files := make(map[string]string)
	fileNames := []string{}
	fileNames = getPossibleFiles(allFiles, files, fileNames)

	log.Println("Done.")

	selected := []string{}
	allKeys := maps.Keys(files)
	selectedKeys := maps.Keys(files)

	if len(allKeys) > 1 {
		selected, selectedKeys = getLanguageSelections(fileNames, selected, files)
	} else {
		selected = []string{files[allKeys[0]]}
	}

	cleanupSelection := getCleanupSelection()
	var outContents string
	if cleanupSelection {
		out := []string{}
		for i, selection := range selected {
			header := langHeader(selectedKeys[i])
			out = append(out, ternary(i > 0, "\n", "")+header+removeUnusedPatterns(selection))
		}
		outContents = strings.Join(out, "\n")
	} else {
		for i, selection := range selected {
			header := langHeader(selectedKeys[i])
			selected[i] = header + selection
		}
		outContents = strings.Join(selected, "\n")
	}

	if fileExists(outFile) {
		handleFileOverwrite(outFile, outContents)
	} else {
		log.Printf("Writing to %s", outFile)
		writeFile(outFile, outContents, true)
	}
}

func langHeader(langName string) string {
	sep := "#========================================================================\n"
	header := fmt.Sprintf(sep+"# %s\n"+sep+"\n", langName)
	return header
}
