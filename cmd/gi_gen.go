package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	. "github.com/chenasraf/gi_gen/internal"
)

func RunMainCmd() {
	helpCommand()

	shouldReturn := cleanCommand()
	if shouldReturn {
		return
	}

	RunGIGen()
}

func cleanCommand() bool {
	var clean *bool = flag.Bool("cache-clean", false, "Clear the .gitignore cache directory, for troubleshooting or for removing trace files of this program.")
	flag.Parse()

	if *clean {
		RemoveCacheDir()
		return true
	}
	return false
}

func helpCommand() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprint(w, "Usage: gi_gen [options]\n\n")
		fmt.Fprint(w, "This program generates .gitignore files for any project. You may clean unused\nlines from the generated file, and the program auto-discovers relevant\n.gitignore templates unless asked not to when the prompt appears.\n\n")
		fmt.Fprint(w, "Run without arguments to use the normal functionality of gi_gen.\n\n")
		fmt.Fprint(w, "Available flags:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nCopyright © 2022 - Chen Asraf\nhttps://casraf.blog\nhttps://github.com/chenasraf/gi_gen\n")
	}

	flag.Bool("help", false, "Display help message")
}

func RunGIGen() {
	wd, err := os.Getwd()
	HandleErr(err)

	outFile := filepath.Join(wd, ".gitignore")
	allFiles, err := PrepareGitignores()
	HandleErr(err)

	fileNames, files := DiscoverRelevantFiles(allFiles)

	log.Println("Done.")

	selected, selectedKeys := GetLanguageSelections(files, fileNames)
	cleanupSelection := AskCleanup()
	outContents := Ternary(cleanupSelection, CleanupMultipleFiles(selected, selectedKeys), GetAllRaw(selected, selectedKeys))

	if FileExists(outFile) {
		HandleFileOverwrite(outFile, outContents)
	} else {
		log.Printf("Writing to %s", outFile)
		WriteFile(outFile, outContents, true)
	}
}
