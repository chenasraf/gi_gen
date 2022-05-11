package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	// "github.com/tcnksm/go-input"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

// create line cache map for skipping previously checked glob lines
var lineCache = make(map[string]bool)

func main() {
	var callerSkip = 0 // might need to be 1 in release mode?
	_, filename, _, _ := runtime.Caller(callerSkip)
	binDir := filepath.Dir(filename)
	sourceDir := filepath.Join(binDir, ".github.gitignore")
	wd, _ := os.Getwd()

	log.Println(binDir)

	if !fileExists(sourceDir) {
		log.Println("Getting gitignore files...")
		runCmd("git", "clone", "--depth=1", repoUrl, sourceDir)
	}

	if getNeedsUpdate() {
		log.Println("Updating gitignore files...")
		runCmd("git", "pull", "origin", "master")
	}

	allFiles, err := getGitignores(sourceDir)
	handleErr(err)

	keep := []string{}

	log.Println("Found:")
	for _, v := range allFiles {
		log.Println("Parsing file: " + v)
		contents := readFile(v)
		lines := strings.Split(contents, "\n")
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
			log.Println("Parsing line: " + line)
			if globExists(filepath.Join(wd, line)) {
				keep = append(keep, line)
			}
		}
	}

	log.Println("Done")

	log.Println("Final output:")
	for _, v := range keep {
		log.Println(v)
	}

}
