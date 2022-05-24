package internal

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

func autoDiscover(allFiles []string) ([]string, map[string]string) {
	answer := AskDiscovery()

	if !answer {
		return allFiles, getAllFiles(allFiles)
	}

	list := discoverByExplicitProjectType()
	if len(list) == 0 {
		list = discoverByExistingPatterns(allFiles)
	}
	return maps.Keys(list), list
}

func getAllFiles(allFiles []string) map[string]string {
	files := make(map[string]string)

	for _, filename := range allFiles {
		contents := ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		files[langName] = contents
	}

	return files
}

func discoverByExistingPatterns(allFiles []string) map[string]string {
	files := make(map[string]string)

	for _, filename := range allFiles {
		contents := ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		if findPatternFileMatches(contents) {
			files[langName] = contents
		}
	}
	return files
}

func discoverByExplicitProjectType() map[string]string {
	wd, err := os.Getwd()
	handleErr(err)

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
		ignoreFile := filepath.Join(getCacheDir(), langName+".gitignore")
		checkFile := filepath.Join(wd, key)

		_, keyExists := results[langName]
		if !keyExists && FileExists(checkFile) {
			results[langName] = ReadFile(ignoreFile)
		}
	}

	return results
}
