package cmd

import (
	"log"
	"os"
	"path/filepath"

	. "github.com/chenasraf/gi_gen/internal"
)

func RunMainCmd() {
	wd, err := os.Getwd()
	HandleErr(err)

	outFile := filepath.Join(wd, ".gitignore")
	allFiles, err := PrepareGitignores()
	HandleErr(err)

	fileNames, files := GetRelevantFiles(allFiles)

	log.Println("Done.")

	selected, selectedKeys := GetLanguageSelections(files, fileNames)
	cleanupSelection := GetCleanupSelection()
	outContents := Ternary(cleanupSelection, CleanupMultiple(selected, selectedKeys), GetAllRaw(selected, selectedKeys))

	if FileExists(outFile) {
		HandleFileOverwrite(outFile, outContents)
	} else {
		log.Printf("Writing to %s", outFile)
		WriteFile(outFile, outContents, true)
	}
}
