package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var repoUrl = "https://github.com/github/gitignore"

func fileExists(path string) bool {
	_, err := os.Stat(path)
	exists := os.IsExist(err)
	return exists
}

func globExists(path string) bool {
	res, err := filepath.Glob(path)
	// log.Println("globExists? " + path)
	handleErr(err)
	return res != nil
}

func getNeedsUpdate() bool {
	localBytes, localErr := exec.Command("git", "rev-parse", "@").Output()
	baseBytes, baseErr := exec.Command("git", "merge-base", "@", "@{u}").Output()
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

func runCmd(cmd string, args ...string) (string, error) {
	res, err := exec.Command(cmd, args...).Output()
	return string(res), err
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
