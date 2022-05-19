package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

var repoUrl = "https://github.com/github/gitignore"

func FileExists(path string) bool {
	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)
	return exists
}

func GlobExists(path string) bool {
	res, err := filepath.Glob(path)
	HandleErr(err)
	return res != nil
}

func GetNeedsUpdate() bool {
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

func RunCmd(cmd string, args ...string) (string, error) {
	res, err := exec.Command(cmd, args...).Output()
	return string(res), err
}

func ReadFile(path string) string {
	res, err := os.ReadFile(path)
	HandleErr(err)
	return string(res)
}

func WriteFile(path string, data string, overwrite bool) bool {
	var err error
	if overwrite {
		// os.Create(path)
		err = os.WriteFile(path, []byte(data), fs.ModeAppend)
		HandleErr(err)

	} else {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		HandleErr(err)
		defer f.Close()
		_, err = f.WriteString("\n" + data)
		HandleErr(err)
	}
	return true
}

func HandleErr(err error) {
	if err != nil {
		log.Println("Encountered an error while running gi_gen:")
		log.Fatalln(err)
	}
}

func Insert[T any](a []T, i int, item T) []T {
	return append(a[:i], append([]T{item}, a[i:]...)...)
}

func Contains[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func Ternary[T any](cond bool, whenTrue T, whenFalse T) T {
	if cond {
		return whenTrue
	}
	return whenFalse
}

func ToString[T any](obj T) string {
	return fmt.Sprint(obj)
}
