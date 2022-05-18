package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

var repoUrl = "https://github.com/github/gitignore"

func fileExists(path string) bool {
	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)
	return exists
}

func globExists(path string) bool {
	res, err := filepath.Glob(path)
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

func readFile(path string) string {
	res, err := os.ReadFile(path)
	handleErr(err)
	return string(res)
}

func writeFile(path string, data string, overwrite bool) bool {
	var err error
	if overwrite {
		// os.Create(path)
		err = os.WriteFile(path, []byte(data), fs.ModeAppend)
	} else {
		f, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		_, err = f.WriteString("\n" + data)
	}
	handleErr(err)
	return true
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func insert[T any](a []T, i int, item T) []T {
	return append(a[:i], append([]T{item}, a[i:]...)...)
}

func contains[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func ternary[T any](cond bool, whenTrue T, whenFalse T) T {
	if cond {
		return whenTrue
	}
	return whenFalse
}
