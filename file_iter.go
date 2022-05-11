package main

import (
	"os"
	"path/filepath"
)

func getGitignores(sourceDir string) ([]string, error) {
	return filepath.Glob(filepath.Join(sourceDir, "*.gitignore"))
}

func readFile(path string) string {
	res, err := os.ReadFile(path)
	handleErr(err)
	return string(res)
}
