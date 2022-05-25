package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chenasraf/gi_gen/internal/utils"
)

func InitCache() ([]string, error) {
	gitignoresDir := GetCacheDir()

	if !utils.FileExists(gitignoresDir) {
		fmt.Println("Getting gitignore files...")
		utils.RunCmd("git", "clone", "--depth=2", utils.RepoUrl, gitignoresDir)
		fmt.Println()
	} else if isCacheNeedsUpdate() {
		fmt.Println("Updating gitignore files...")
		utils.RunCmd("git", "-C", gitignoresDir, "pull", "origin", "main")
		fmt.Println()
	}

	return getGitignoreFiles(gitignoresDir)
}

func RemoveCacheDir() {
	cacheDir := GetCacheDir()
	fmt.Printf("Removing cache directory: %s...\n", cacheDir)
	os.RemoveAll(cacheDir)
	fmt.Println("Done")
}

func GetCacheDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".github.gitignore")
}
