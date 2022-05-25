package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/chenasraf/gi_gen/internal"
)

func RunMainCmd() {
	initFlags()
	initHelpCommand()
	shouldReturn := false

	flag.Parse()

	shouldReturn = cleanCommand()
	if shouldReturn {
		return
	}

	flagLangs := getLangsFromArgs()
	internal.GIGen(&internal.GIGenOptions{
		Languages:     &flagLangs,
		CleanOutput:   &cleanOutput,
		OverwriteFile: &overwriteFile,
		AppendFile:    &appendFile,
	})
}

var langsRaw string = ""
var cleanCache bool = false
var cleanOutput bool
var overwriteFile bool
var appendFile bool

func shorthand(msg string) string {
	return msg + " (shorthand)"
}

func initFlags() {
	appendUsage := "Append to .gitignore file if it already exists"
	langsUsage := "List the languages you want to use as templates.\nTo add multiple templates, use commas as separators, e.g.: -languages Node,Python"
	clearCacheUsage := "Clear the .gitignore cache directory, for troubleshooting or for removing trace files of this program. Exits after running, so other flags will be ignored."
	cleanOutputUsage := "Perform cleanup on the output .gitignore file, removing any unused patterns"
	overwriteUsage := "Overwrite .gitignore file if it already exists"

	flag.Bool("help", false, "Display help message")
	flag.BoolVar(&cleanCache, "clear-cache", false, clearCacheUsage)
	flag.BoolVar(&cleanOutput, "c", false, shorthand(cleanOutputUsage))
	flag.BoolVar(&cleanOutput, "clean-output", false, cleanOutputUsage)
	flag.BoolVar(&overwriteFile, "w", false, shorthand(overwriteUsage))
	flag.BoolVar(&overwriteFile, "overwrite", false, overwriteUsage)
	flag.BoolVar(&appendFile, "a", false, shorthand(appendUsage))
	flag.BoolVar(&appendFile, "append", false, appendUsage)
	flag.StringVar(&langsRaw, "l", langsRaw, shorthand(langsUsage))
	flag.StringVar(&langsRaw, "languages", langsRaw, langsUsage)
}

func getLangsFromArgs() []string {
	return strings.Split(langsRaw, ",")
}

func cleanCommand() bool {
	if cleanCache {
		internal.RemoveCacheDir()
		return true
	}
	return false
}

func initHelpCommand() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprint(w, "Usage: gi_gen [options]\n\n")
		fmt.Fprint(w, "This program generates .gitignore files for any project. You may clean unused\nlines from the generated file, and the program auto-discovers relevant\n.gitignore templates unless asked not to when the prompt appears.\n\n")
		fmt.Fprint(w, "Run without arguments to use the normal functionality of gi_gen.\n\n")
		fmt.Fprint(w, "Available flags:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nCopyright © 2022 - Chen Asraf\nhttps://casraf.blog\nhttps://github.com/chenasraf/gi_gen\n")
	}
}
