package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func initCache() ([]string, error) {
	gitignoresDir := getCacheDir()

	if !fileExists(gitignoresDir) {
		fmt.Println("Getting gitignore files...")
		runCmd("git", "clone", "--depth=2", repoUrl, gitignoresDir)
		fmt.Println()
	} else if isCacheNeedsUpdate() {
		fmt.Println("Updating gitignore files...")
		runCmd("git", "-C", gitignoresDir, "pull", "origin", "main")
		fmt.Println()
	}

	return getGitignoreFiles(gitignoresDir)
}

func RemoveCacheDir() {
	cacheDir := getCacheDir()
	fmt.Printf("Removing cache directory: %s...\n", cacheDir)
	os.RemoveAll(cacheDir)
	fmt.Println("Done")
}

func getCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".github.gitignore")
}