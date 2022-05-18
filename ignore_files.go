package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func prepareGitignores() ([]string, error) {
	gitignoresDir := getCacheDir()

	if !fileExists(gitignoresDir) {
		log.Println("Getting gitignore files...")
		runCmd("git", "clone", "--depth=1", repoUrl, gitignoresDir)
	}

	if getNeedsUpdate() {
		log.Println("Updating gitignore files...")
		runCmd("git", "pull", "origin", "master")
		os.RemoveAll(filepath.Join(gitignoresDir, ".git"))
	}

	return getGitignores(gitignoresDir)
}

func getCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	gitignoresDir := filepath.Join(homeDir, ".github.gitignore")
	return gitignoresDir
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
		line = strings.TrimSpace(line)

		// ignore empty lines / comments
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

func removeUnusedPatterns(contents string) string {
	wd, _ := os.Getwd()
	lines := strings.Split(contents, "\n")
	keep := []string{}
	lastTakenIdx := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if len(trimmed) == 0 || trimmed[0] == '#' {
			continue
			// keep = append(keep, line)
			// continue
		}

		if globExists(filepath.Join(wd, line)) {
			if i > 0 {
				j := 1
				foundComment := false
				comments := []string{}
				for {
					if i-j < 0 || i-j <= lastTakenIdx {
						break
					}
					cur := lines[i-j]
					if len(cur) > 0 && cur[0] != '#' {
						if !foundComment {
						} else {
							break
						}
					} else {
						lastTakenIdx = i - j
						if len(cur) > 0 && cur[0] == '#' {
							foundComment = true
						}
						comments = insert(comments, 0, cur)
					}
					j++
				}
				for _, v := range comments {
					keep = append(keep, v)
				}
			}
			keep = append(keep, line)
		}
	}

	return strings.Join(keep, "\n")
}
