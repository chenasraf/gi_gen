package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chenasraf/gi_gen/internal/utils"
	"golang.org/x/exp/maps"
)

type GIGenOptions struct {
	Languages         *[]string
	CleanOutput       bool
	CleanOutputUsed   bool
	KeepOutput        bool
	KeepOutputUsed    bool
	AutoDiscover      bool
	AutoDiscoverUsed  bool
	OverwriteFile     bool
	OverwriteFileUsed bool
	AppendFile        bool
	AppendFileUsed    bool
}

func GIGen(options *GIGenOptions) {
	wd, err := os.Getwd()
	utils.HandleErr(err)
	opts := utils.Ternary(options != nil, *options, GIGenOptions{})

	outFile := filepath.Join(wd, ".gitignore")
	allFiles, err := InitCache()
	cacheDir := GetCacheDir()
	utils.HandleErr(err)
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

	if utils.FileExists(outFile) {
		overwriteSelection := getOverwriteSelection(opts)
		utils.HandleFileOverwrite(outFile, outContents, overwriteSelection)
	} else {
		fmt.Println()
		fmt.Printf("Writing to %s\n", outFile)
		utils.WriteFile(outFile, outContents, true)
	}

	fmt.Println()
	fmt.Println("Done.")
}

func getOverwriteSelection(opts GIGenOptions) string {
	var overwriteSelection string
	if opts.OverwriteFileUsed || opts.AppendFileUsed {
		overwriteSelection = utils.Ternary(opts.OverwriteFileUsed, "Overwrite", "Append")
	} else {
		return askOverwrite()
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
	if opts.CleanOutputUsed || opts.KeepOutput {
		cleanupSelection = opts.CleanOutput && !opts.KeepOutput
	} else {
		cleanupSelection = askCleanup()
	}
	return cleanupSelection
}

func getProcessFiles(
	opts GIGenOptions, cacheDir string, fileNames []string, files map[string]string, allFiles []string,
) ([]string, []string, map[string]string) {
	mappedFileNames := []string{}

	if len(*opts.Languages) > 0 && (*opts.Languages)[0] != "" {
		for _, lng := range *opts.Languages {
			filePath := filepath.Join(cacheDir, lng+".gitignore")
			if utils.FileExists(filePath) {
				mappedFileNames = append(mappedFileNames, filePath)
			}
		}
		fileNames, files = mappedFileNames, getAllFiles(mappedFileNames)
	} else {
		fileNames, files = readFromSelections(allFiles, opts)
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
