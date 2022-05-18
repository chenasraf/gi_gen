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

	fileNames, files := getPossibleFiles(allFiles)

	log.Println("Done.")

	selected, selectedKeys := getLanguages(files, fileNames)

	cleanupSelection := getCleanupSelection()
	var outContents string
	if cleanupSelection {
		out := []string{}
		for i, selection := range selected {
			cleanSelection := removeUnusedPatterns(selection)
			if strings.TrimSpace(cleanSelection) == "" {
				continue
			}
			header := langHeader(selectedKeys[i])
			prefixNewline := ternary(i > 0, "\n", "")
			contents := prefixNewline + header + cleanSelection
			out = append(out, contents)
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

func getLanguages(files map[string]string, fileNames []string) ([]string, []string) {
	selected := []string{}
	allKeys := maps.Keys(files)
	selectedKeys := maps.Keys(files)

	if len(allKeys) > 1 {
		selected, selectedKeys = getLanguageSelections(fileNames, selected, files)
	} else {
		selected = []string{files[allKeys[0]]}
	}

	return selected, selectedKeys
}

func langHeader(langName string) string {
	sep := "#========================================================================\n"
	header := fmt.Sprintf(sep+"# %s\n"+sep+"\n", langName)
	return header
}
