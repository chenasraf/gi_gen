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

	mappedFileNames, fileNames, files := getProcessFiles(opts, cacheDir, fileNames, files, allFiles)
	selectedContents, selectedKeys := getSelections(mappedFileNames, files, fileNames)

	if len(mappedFileNames) == 0 {
		fmt.Println()
	}
	fmt.Printf("Selected languages: %s\n", strings.Join(selectedKeys, ", "))

	cleanupSelection := getCleanupSelection(opts)
	outContents := processFileOutput(cleanupSelection, selectedContents, selectedKeys)

	if fileExists(outFile) {
		getOverwriteSelection := newFunction(opts)
		handleFileOverwrite(outFile, outContents, getOverwriteSelection)
	} else {
		fmt.Println()
		fmt.Printf("Writing to %s\n", outFile)
		writeFile(outFile, outContents, true)
	}

	fmt.Println()
	fmt.Println("Done.")
}

func newFunction(opts GIGenOptions) string {
	var overwriteSelection string
	if opts.OverwriteFile != nil || opts.AppendFile != nil {
		overwriteSelection = ternary(opts.OverwriteFile != nil, "Overwrite", "Append")
	} else {
		askOverwrite()
	}
	return overwriteSelection
}

func processFileOutput(cleanupSelection bool, selectedContents []string, selectedKeys []string) string {
	var outContents string
	if cleanupSelection {
		outContents = cleanupMultipleFiles(selectedContents, selectedKeys)
	} else {
		outContents = getAllRaw(selectedContents, selectedKeys)
	}
	return outContents
}

func getCleanupSelection(opts GIGenOptions) bool {
	var cleanupSelection bool
	if opts.CleanOutput != nil {
		cleanupSelection = *opts.CleanOutput
	} else {
		cleanupSelection = askCleanup()
	}
	return cleanupSelection
}

func getProcessFiles(
	opts GIGenOptions, cacheDir string, fileNames []string, files map[string]string, allFiles []string,
) ([]string, []string, map[string]string) {
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
	return mappedFileNames, fileNames, files
}

func getSelections(mappedFileNames []string, files map[string]string, fileNames []string) ([]string, []string) {
	var selectedContents []string
	var selectedKeys []string
	if len(mappedFileNames) > 0 {
		selectedContents, selectedKeys = maps.Values(files), maps.Keys(files)
	} else {
		selectedContents, selectedKeys = getLanguages(files, fileNames)
	}
	return selectedContents, selectedKeys
}
