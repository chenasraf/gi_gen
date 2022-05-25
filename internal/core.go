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
	allFiles, err := initCache()
	handleErr(err)

	fileNames, files := autoDiscover(allFiles)

	selected, selectedKeys := getLanguages(files, fileNames)
	cleanupSelection := askCleanup()
	outContents := ternary(cleanupSelection, cleanupMultipleFiles(selected, selectedKeys), getAllRaw(selected, selectedKeys))

	if fileExists(outFile) {
		handleFileOverwrite(outFile, outContents)
	} else {
		fmt.Println()
		fmt.Printf("Writing to %s\n", outFile)
		writeFile(outFile, outContents, true)
	}

	fmt.Println()
	fmt.Println("Done.")
}
