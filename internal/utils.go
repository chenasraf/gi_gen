package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
		err = os.WriteFile(path, []byte(data), 0644)
		handleErr(err)
	} else {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		handleErr(err)
		defer f.Close()
		_, err = f.WriteString("\n" + data)
		handleErr(err)
	}
	return true
}

func handleErr(err error) {
	if err != nil {
		fmt.Println("Encountered an error while running gi_gen:")
		fmt.Println(err)
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

func toString[T any](obj T) string {
	return fmt.Sprint(obj)
}

func handleFileOverwrite(path string, contents string, selection string) {
	switch selection {
	case "Skip":
		quit("Nothing to do, exiting")
		return
	case "Overwrite":
		fmt.Println()
		fmt.Printf("Writing to %s\n", path)
		writeFile(path, contents, true)
		return
	case "Append":
		fmt.Println()
		fmt.Printf("Appending to %s\n", path)
		writeFile(path, contents, false)
		return
	}
	quit("Bad selection")
}

func keyInterrupt() {
	quit("KeyInterrupt: Quitting")
}

func quit(message string) {
	fmt.Println()
	fmt.Println(message)
	os.Exit(1)
}
