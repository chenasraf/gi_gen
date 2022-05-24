package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

func DiscoverRelevantFiles(allFiles []string) ([]string, map[string]string) {
	answer := AskDiscovery()

	if !answer {
		return allFiles, getAll(allFiles)
	}

	byType := discoverProjectType()
	return maps.Keys(byType), byType
}

func getAll(allFiles []string) map[string]string {
	files := make(map[string]string)

	for _, filename := range allFiles {
		contents := ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		files[langName] = contents
	}

	return files
}

func discoverExistingPatterns(allFiles []string, answer bool, files map[string]string) map[string]string {
	for _, filename := range allFiles {
		contents := ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		if answer {
			if FindFileMatches(contents) {
				files[langName] = contents
			}
		} else {
			files[langName] = contents
		}
	}
	return files
}

func discoverProjectType() map[string]string {
	wd, err := os.Getwd()
	HandleErr(err)

	discoveryMap := make(map[string]string)
	discoveryMap["package.json"] = "Node"
	discoveryMap["tsconfig.json"] = "Node"
	discoveryMap["jsconfig.json"] = "Node"
	discoveryMap["setup.py"] = "Python"
	discoveryMap["__init__.py"] = "Python"
	discoveryMap["lib/main.dart"] = "Dart"
	discoveryMap["pubspec.yaml"] = "Dart"
	discoveryMap["pubspec.yml"] = "Dart"
	discoveryMap["go.mod"] = "Go"
	discoveryMap["go.sum"] = "Go"
	discoveryMap["main.go"] = "Go"

	results := make(map[string]string)

	for _, key := range maps.Keys(discoveryMap) {
		langName := discoveryMap[key]
		ignoreFile := filepath.Join(GetCacheDir(), langName+".gitignore")
		checkFile := filepath.Join(wd, key)

		fmt.Println("Trying file " + checkFile)
		_, keyExists := results[langName]
		if !keyExists && FileExists(checkFile) {
			fmt.Println("Found lang " + langName)
			results[langName] = ReadFile(ignoreFile)
		}
	}

	return results
}
