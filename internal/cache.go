package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	return getCacheTemplates(gitignoresDir)
}

func getCacheTemplates(sourceDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(sourceDir, "*.gitignore"))
}

func isCacheNeedsUpdate() bool {
	gitignoresDir := GetCacheDir()
	localBytes, localErr := exec.Command("git", "-C", gitignoresDir, "rev-list", "--count", "HEAD..@{u}").Output()
	utils.HandleErr(localErr)
	localStr := strings.TrimSpace(string(localBytes))

	return localStr != "0"
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
