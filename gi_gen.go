package main

import (
	"fmt"
	"os/exec"
)

// UNUSED allows unused variables to be included in Go programs
func UNUSED(x ...interface{}) {}

func main() {
	var repoUrl = "https://github.com/github/gitignore"
	var cmd = exec.Command("git", "clone", "--no-checkout", "--depth=1", repoUrl, ".source")
	cmd.Run()
	fmt.Println("Done")
}