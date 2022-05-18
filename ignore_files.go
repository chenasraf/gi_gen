package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func prepareGitignores() ([]string, error) {
	var callerSkip = 0 // might need to be 1 in release mode?
	_, filename, _, _ := runtime.Caller(callerSkip)
	binDir := filepath.Dir(filename)
	gitignoresDir := filepath.Join(binDir, ".github.gitignore")

	if !fileExists(gitignoresDir) {
		log.Println("Getting gitignore files...")
		runCmd("git", "clone", "--depth=1", repoUrl, gitignoresDir)
	}

	if getNeedsUpdate() {
		log.Println("Updating gitignore files...")
		runCmd("git", "pull", "origin", "master")
	}

	return getGitignores(gitignoresDir)
}

func getGitignores(sourceDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(sourceDir, "*.gitignore"))
}

var ignoreLines = []string{
	"/*",
	".",
	".vscode",
	".vscode/*",
	".idea",
	".idea/*",
}

func findFileMatches(patterns string) bool {
	lines := strings.Split(patterns, "\n")
	wd, _ := os.Getwd()

	for _, line := range lines {
		// ignore empty lines / comments
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.ToLower(line)[0] == '#' {
			continue
		}
		idx := strings.Index(line, "#")
		// ignore comments at end of line
		if idx > -1 && (idx == 0 || line[idx-1] != '\\') {
			line = strings.TrimSpace(line[0:idx])
		}
		if len(line) == 0 || contains(ignoreLines, line) {
			continue
		}
		if globExists(filepath.Join(wd, line)) {
			return true
		}
	}

	return false
}

func removeUnusedPatterns(filename string, keep []string) []string {
	wd, _ := os.Getwd()
	contents := readFile(filename)
	lines := strings.Split(contents, "\n")
	keepCopy := []string{}
	copy(keep, keepCopy)
	found := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		_, exists := lineCache[trimmed]
		if exists {
			continue
		}
		lineCache[trimmed] = true
		if len(trimmed) == 0 || strings.ToLower(string(trimmed[0])) == "#" {
			lineCache[trimmed] = false
			continue
		}

		if globExists(filepath.Join(wd, line)) {
			found = true
			keepCopy = append(keepCopy, line)
		}
	}

	if found {
		keepCopy = insert(keepCopy, 0, fmt.Sprintf("\n# %s", filepath.Base(filename)))
	}

	return keep
}
