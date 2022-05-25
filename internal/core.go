package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

type GIGenOptions struct {
	Languages     *[]string
	CleanOutput   *bool
	OverwriteFile *bool
	AppendFile    *bool
}

func GIGen(options *GIGenOptions) {
	wd, err := os.Getwd()
	handleErr(err)
	opts := ternary(options != nil, *options, GIGenOptions{})

	outFile := filepath.Join(wd, ".gitignore")
	allFiles, err := InitCache()
	cacheDir := GetCacheDir()
	handleErr(err)
	var fileNames []string
	var files map[string]string

	mappedFileNames := []string{}

	if len(*opts.Languages) > 0 {
		for _, lng := range *opts.Languages {
			filePath := filepath.Join(cacheDir, lng+".gitignore")
			if fileExists(filePath) {
				mappedFileNames = append(mappedFileNames, filePath)
			}
		}
		fileNames, files = mappedFileNames, getAllFiles(mappedFileNames)
	} else {
		fileNames, files = autoDiscover(allFiles)
	}

	var selectedContents []string
	var selectedKeys []string
	if len(mappedFileNames) > 0 {
		selectedContents, selectedKeys = maps.Values(files), maps.Keys(files)
	} else {
		selectedContents, selectedKeys = getLanguages(files, fileNames)
	}
	if len(mappedFileNames) == 0 {
		fmt.Println()
	}
	fmt.Printf("Selected languages: %s\n", strings.Join(selectedKeys, ", "))

	var cleanupSelection bool
	if opts.CleanOutput != nil {
		cleanupSelection = *opts.CleanOutput
	} else {
		cleanupSelection = askCleanup()
	}

	var outContents string
	if cleanupSelection {
		outContents = cleanupMultipleFiles(selectedContents, selectedKeys)
	} else {
		outContents = getAllRaw(selectedContents, selectedKeys)
	}

	if fileExists(outFile) {
		var overwriteSelection string
		if opts.OverwriteFile != nil || opts.AppendFile != nil {
			overwriteSelection = ternary(opts.OverwriteFile != nil, "Overwrite", "Append")
		} else {
			askOverwrite()
		}
		handleFileOverwrite(outFile, outContents, overwriteSelection)
	} else {
		fmt.Println()
		fmt.Printf("Writing to %s\n", outFile)
		writeFile(outFile, outContents, true)
	}

	fmt.Println()
	fmt.Println("Done.")
}
