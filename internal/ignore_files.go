package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

func PrepareGitignores() ([]string, error) {
	gitignoresDir := GetCacheDir()

	if !FileExists(gitignoresDir) {
		log.Println("Getting gitignore files...")
		RunCmd("git", "clone", "--depth=1", repoUrl, gitignoresDir)
	} else if GetNeedsUpdate() {
		log.Println("Updating gitignore files...")
		RunCmd("git", "-C", gitignoresDir, "pull", "origin", "master")
	}

	return GetGitignores(gitignoresDir)
}

func RemoveCacheDir() {
	cacheDir := GetCacheDir()
	log.Printf("Removing cache directory: %s...\n", cacheDir)
	os.RemoveAll(cacheDir)
	log.Println("Done")
}

func GetCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".github.gitignore")
}

func GetGitignores(sourceDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(sourceDir, "*.gitignore"))
}

func GetNeedsUpdate() bool {
	gitignoresDir := GetCacheDir()
	localBytes, localErr := exec.Command("git", "-C", gitignoresDir, "rev-parse", "@").Output()
	baseBytes, baseErr := exec.Command("git", "-C", gitignoresDir, "merge-base", "@", "@{u}").Output()
	if localErr != nil {
		log.Fatal(localErr)
		os.Exit(1)
	}
	if baseErr != nil {
		log.Fatal(baseErr)
		os.Exit(1)
	}
	localStr := string(localBytes)
	baseStr := string(baseBytes)

	return localStr == baseStr
}

var ignoreLines = []string{
	"/*",
	".",
	".vscode",
	".vscode/*",
	".idea",
	".idea/*",
}

func FindFileMatches(patterns string) bool {
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

func RemoveUnusedPatterns(contents string) string {
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
				keep = GatherPreviousCommentGroup(i, lastTakenIdx, lines, keep)
			}

			keep = append(keep, line)
		}
	}

	return strings.Join(keep, "\n")
}

func GatherPreviousCommentGroup(i int, lastTakenIdx int, lines []string, keep []string) []string {
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

func GetLanguageSelections(files map[string]string, fileNames []string) ([]string, []string) {
	selected := []string{}
	allKeys := maps.Keys(files)
	selectedKeys := maps.Keys(files)

	if len(allKeys) == 0 {
		selected = []string{}
	} else if len(allKeys) > 1 {
		selected, selectedKeys = AskLanguage(fileNames, selected, files)
	} else {
		selected = []string{files[allKeys[0]]}
	}

	return selected, selectedKeys
}

func LangHeader(langName string) string {
	sep := "#========================================================================\n"
	header := fmt.Sprintf(sep+"# %s\n"+sep+"\n", langName)
	return header
}

func GetAllRaw(selected []string, selectedKeys []string) string {
	for i, selection := range selected {
		header := Ternary(len(selected) > 1, LangHeader(selectedKeys[i]), "")
		selected[i] = header + selection
	}
	return strings.Join(selected, "\n")
}

func CleanupMultipleFiles(files []string, langKeys []string) string {
	out := []string{}
	for i, selection := range files {
		cleanSelection := RemoveUnusedPatterns(selection)
		if strings.TrimSpace(cleanSelection) == "" {
			continue
		}
		header := Ternary(len(files) > 1, LangHeader(langKeys[i]), "")
		prefixNewline := Ternary(i > 0, "\n", "")
		contents := prefixNewline + header + cleanSelection
		out = append(out, contents)
	}
	return strings.Join(out, "\n")
}
