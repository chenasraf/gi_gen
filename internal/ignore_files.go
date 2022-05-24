package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

func getGitignoreFiles(sourceDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(sourceDir, "*.gitignore"))
}

func isCacheNeedsUpdate() bool {
	gitignoresDir := getCacheDir()
	localBytes, localErr := exec.Command("git", "-C", gitignoresDir, "rev-list", "--count", "HEAD..@{u}").Output()
	handleErr(localErr)
	localStr := strings.TrimSpace(string(localBytes))

	return localStr != "0"
}

var ignoreLines = []string{
	"/*",
	".",
	".vscode",
	".vscode/*",
	".idea",
	".idea/*",
}

func findPatternFileMatches(patterns string) bool {
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

		if len(line) == 0 || Contains(ignoreLines, line) {
			continue
		}
		if GlobExists(filepath.Join(wd, line)) {
			return true
		}
	}

	return false
}

var patternCache []string = []string{}

func removeUnusedPatterns(contents string) string {
	wd, _ := os.Getwd()
	lines := strings.Split(contents, "\n")
	keep := []string{}
	lastTakenIdx := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if len(trimmed) == 0 || trimmed[0] == '#' {
			continue
		}

		if GlobExists(filepath.Join(wd, trimmed)) {
			if Contains(patternCache, trimmed) {
				continue
			}

			patternCache = append(patternCache, trimmed)

			if i > 0 {
				keep = gatherPreviousCommentGroup(i, lastTakenIdx, lines, keep)
			}

			keep = append(keep, line)
		}
	}

	return strings.Join(keep, "\n")
}

func gatherPreviousCommentGroup(i int, lastTakenIdx int, lines []string, keep []string) []string {
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
			comments = Insert(comments, 0, cur)
		}
		j++
	}
	for _, v := range comments {
		keep = append(keep, v)
	}
	return keep
}

func getLanguages(files map[string]string, fileNames []string) ([]string, []string) {
	selected := []string{}
	allKeys := maps.Keys(files)
	selectedKeys := maps.Keys(files)
	fmt.Println()
	if len(allKeys) == 0 {
		fmt.Println("Found no templates. Quitting.")
		os.Exit(1)
	} else if len(allKeys) > 1 {
		fmt.Println("Found " + fmt.Sprint(len(fileNames)) +
			" possible matches in your project for gitignore files.")
		selected, selectedKeys = askLanguage(fileNames, selected, files)
	} else {
		fmt.Printf("Found one match for your project: %s. Proceeding...\n", allKeys[0])
		selected = []string{files[allKeys[0]]}
	}

	return selected, selectedKeys
}

func langHeader(langName string) string {
	sep := "#========================================================================\n"
	header := fmt.Sprintf(sep+"# %s\n"+sep+"\n", langName)
	return header
}

func getAllRaw(selected []string, selectedKeys []string) string {
	for i, selection := range selected {
		header := Ternary(len(selected) > 1, langHeader(selectedKeys[i]), "")
		selected[i] = header + selection
	}
	return strings.Join(selected, "\n")
}

func cleanupMultipleFiles(files []string, langKeys []string) string {
	out := []string{}
	for i, selection := range files {
		cleanSelection := removeUnusedPatterns(selection)
		if strings.TrimSpace(cleanSelection) == "" {
			continue
		}
		header := Ternary(len(files) > 1, langHeader(langKeys[i]), "")
		prefixNewline := Ternary(i > 0, "\n", "")
		contents := prefixNewline + header + cleanSelection
		out = append(out, contents)
	}
	return strings.Join(out, "\n")
}
