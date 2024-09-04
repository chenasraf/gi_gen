package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chenasraf/utils"
)

var ignoreLines = []string{
	"/*",
	".",
	".vscode",
	".vscode/*",
	".idea",
	".idea/*",
}

func findPatternFileMatches(patterns string) (bool, string) {
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

		if len(line) == 0 || utils.SliceContains(ignoreLines, line) {
			continue
		}
		if utils.GlobExists(filepath.Join(wd, line)) {
			return true, line
		}
	}

	return false, ""
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

		if utils.GlobExists(filepath.Join(wd, trimmed)) {
			if utils.SliceContains(patternCache, trimmed) {
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
			comments = utils.Insert(comments, 0, cur)
		}
		j++
	}
	for _, v := range comments {
		keep = append(keep, v)
	}
	return keep
}

func langHeader(langName string) string {
	sep := "#========================================================================\n"
	header := fmt.Sprintf(sep+"# %s\n"+sep+"\n", langName)
	return header
}

func getAllRaw(selected []string, selectedKeys []string) string {
	for i, selection := range selected {
		header := utils.Ternary(len(selected) > 1, langHeader(selectedKeys[i]), "")
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
		header := utils.Ternary(len(files) > 1, langHeader(langKeys[i]), "")
		prefixNewline := utils.Ternary(i > 0, "\n", "")
		contents := prefixNewline + header + cleanSelection
		// exclude dupes
		fmt.Println("cleanSelection: " + cleanSelection)
		if !utils.SliceContains(out, cleanSelection) {
			out = append(out, contents)
		} else {
			fmt.Println("Skipped duplicate: " + cleanSelection)
		}
	}
	return strings.Join(out, "\n")
}
