package cmd

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/chenasraf/gi_gen/internal"
	"github.com/chenasraf/gi_gen/internal/utils"
)

func RunMainCmd() {
	initFlags()
	initHelpCommand()
	shouldReturn := false

	flag.Parse()

	shouldReturn = allLanguageCommand()
	if shouldReturn {
		return
	}

	shouldReturn = detectLanguageCommand()
	if shouldReturn {
		return
	}

	shouldReturn = cleanCommand()
	if shouldReturn {
		return
	}

	flagLangs := getLangsFromArgs()
	internal.GIGen(&internal.GIGenOptions{
		Languages:         &flagLangs,
		AutoDiscover:      autoDiscover,
		AutoDiscoverUsed:  isFlagPassed("auto-discover") || isFlagPassed("d"),
		CleanOutput:       cleanOutput,
		CleanOutputUsed:   isFlagPassed("clean-output") || isFlagPassed("c"),
		KeepOutput:        keepOutput,
		KeepOutputUsed:    isFlagPassed("keep-output") || isFlagPassed("k"),
		OverwriteFile:     overwriteFile,
		OverwriteFileUsed: isFlagPassed("overwrite") || isFlagPassed("w"),
		AppendFile:        appendFile,
		AppendFileUsed:    isFlagPassed("append") || isFlagPassed("a"),
	})
}

var langsRaw string = ""
var cleanCache bool
var cleanOutput bool
var keepOutput bool
var overwriteFile bool
var appendFile bool
var detectLanguages bool
var allLanguages bool
var autoDiscover bool

func shorthand(msg string) string {
	return msg + " (shorthand)"
}

func initFlags() {
	langsUsage := "List the languages you want to use as templates.\n" +
		"To add multiple templates, use commas as separators, e.g.: -languages Node,Python"
	autoDiscoverUsage := "Use auto-discovery for project, detecting the project type and using the result as the pre-" +
		"selected template list."
	cleanOutputUsage := "Perform cleanup on the output .gitignore file, removing any unused patterns"
	keepOutputUsage := "Do not perform cleanup on the output .gitignore file, keep all the original contents " +
		"(opposite of -clean-output)"
	appendUsage := "Append to .gitignore file if it already exists"
	overwriteUsage := "Overwrite .gitignore file if it already exists"
	clearCacheUsage := "Clear the .gitignore cache directory, for troubleshooting or for removing trace files of this " +
		"program. Exits after running, so other flags will be ignored."
	detectLanguagesUsage := "Outputs the automatically-detected languages, separated by newlines, and exits. Useful " +
		"for outside tools detection."
	allLanguagesUsage := "Outputs all the available languages, separated by newlines, and exits. Useful for " +
		"outside tools detection."

	flag.Bool("help", false, "Display help message")
	flag.BoolVar(&cleanCache, "clear-cache", false, clearCacheUsage)

	flag.BoolVar(&autoDiscover, "d", false, shorthand(autoDiscoverUsage))
	flag.BoolVar(&autoDiscover, "auto-discover", false, autoDiscoverUsage)

	flag.BoolVar(&cleanOutput, "c", false, shorthand(cleanOutputUsage))
	flag.BoolVar(&cleanOutput, "clean-output", false, cleanOutputUsage)

	flag.BoolVar(&keepOutput, "k", false, shorthand(keepOutputUsage))
	flag.BoolVar(&keepOutput, "keep-output", false, keepOutputUsage)

	flag.BoolVar(&overwriteFile, "w", false, shorthand(overwriteUsage))
	flag.BoolVar(&overwriteFile, "overwrite", false, overwriteUsage)

	flag.BoolVar(&appendFile, "a", false, shorthand(appendUsage))
	flag.BoolVar(&appendFile, "append", false, appendUsage)

	flag.BoolVar(&detectLanguages, "detect-languages", false, detectLanguagesUsage)
	flag.BoolVar(&allLanguages, "all-languages", false, allLanguagesUsage)

	flag.StringVar(&langsRaw, "l", langsRaw, shorthand(langsUsage))
	flag.StringVar(&langsRaw, "languages", langsRaw, langsUsage)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
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

func detectLanguageCommand() bool {
	if detectLanguages {
		allFiles, err := internal.InitCache()
		discovery, _ := internal.AutoDiscover(allFiles)
		utils.HandleErr(err)
		fmt.Println(strings.Join(discovery, "\n"))
		return true
	}
	return false
}

func allLanguageCommand() bool {
	if allLanguages {
		allFiles, err := internal.InitCache()
		utils.HandleErr(err)
		out := []string{}
		for _, fn := range allFiles {
			basename := filepath.Base(fn)
			langName := basename[:strings.Index(basename, ".")]
			out = append(out, langName)
		}
		fmt.Println(strings.Join(out, "\n"))
		return true
	}
	return false
}

func initHelpCommand() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprint(w, "Usage: gi_gen [options]\n\n")
		fmt.Fprint(w, "This program generates .gitignore files for any project. You may clean unused\n"+
			"lines from the generated file, and the program auto-discovers relevant\n."+
			"gitignore templates unless asked not to when the prompt appears.\n\n")
		fmt.Fprint(w, "Run without arguments to use the normal functionality of gi_gen.\n\n")
		fmt.Fprint(w, "Available flags:\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(w, "\nCopyright © 2022 - Chen Asraf\nhttps://casraf.blog\nhttps://github.com/chenasraf/gi_gen\n")
	}
}
