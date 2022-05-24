package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func GIGen() {
	wd, err := os.Getwd()
	handleErr(err)

	outFile := filepath.Join(wd, ".gitignore")
	allFiles, err := prepareGitignores()
	handleErr(err)

	fileNames, files := autoDiscover(allFiles)

	selected, selectedKeys := getLanguageSelections(files, fileNames)
	cleanupSelection := AskCleanup()
	outContents := Ternary(cleanupSelection, cleanupMultipleFiles(selected, selectedKeys), getAllRaw(selected, selectedKeys))

	if FileExists(outFile) {
		HandleFileOverwrite(outFile, outContents)
	} else {
		fmt.Println()
		fmt.Printf("Writing to %s\n", outFile)
		WriteFile(outFile, outContents, true)
	}

	fmt.Println()
	fmt.Println("Done.")
}
