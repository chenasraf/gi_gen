package main

import (
	"os"
	"strings"

	"github.com/Cleverse/go-utilities/nullable"
	"github.com/chenasraf/goutils"
)

// OutputBehavior
//
// - clean   (0): remove ignore patterns not matching actual files
//
// - keep    (1): keep ignore patterns even if they don't match actual files
type OutputBehavior int

const (
	BehaviorClean OutputBehavior = iota
	BehaviorKeep
)

// ConflictBehavior
//
// - skip       (0): skip the gitignore file if it already exists
//
// - merge      (1): merge the new gitignore file with the existing one if it exists
//
// - overwrite  (2): overwrite the existing gitignore file with the new one
type ConflictBehavior int

const (
	BehaviorSkip ConflictBehavior = iota
	BehaviorMerge
	BehaviorOverwrite
)

// main options type
type Options struct {
	languages        *[]string
	outputBehavior   nullable.Nullable[OutputBehavior]
	autoSelect       bool
	conflictBehavior nullable.Nullable[ConflictBehavior]
}

func parseArgs() *Options {
	opts := &Options{
		languages:        &[]string{},
		outputBehavior:   nullable.New[OutputBehavior](),
		autoSelect:       false,
		conflictBehavior: nullable.New[ConflictBehavior](),
	}

	args := os.Args[1:]

	for len(args) > 0 {
		switch args[0] {
		case "--languages", "-l":
			{
				split := strings.Split(args[1], ",")
				*opts.languages = utils.MapSlice(split, strings.TrimSpace)
				opts.autoSelect = true
				args = args[2:]
			}
		case "--auto-select", "-a":
			{
				opts.autoSelect = true
				args = args[1:]
			}
		case "--clean", "-c":
			{
				opts.outputBehavior.Set(BehaviorClean)
				args = args[1:]
			}
		case "--keep", "-k":
			{
				opts.outputBehavior.Set(BehaviorKeep)
				args = args[1:]
			}
		case "--skip", "-s":
			{
				opts.conflictBehavior.Set(BehaviorSkip)
				args = args[1:]
			}
		case "--merge", "-m":
			{
				opts.conflictBehavior.Set(BehaviorMerge)
				args = args[1:]
			}
		case "--overwrite", "-w":
			{
				opts.conflictBehavior.Set(BehaviorOverwrite)
				args = args[1:]
			}
		default:
			args = args[1:]
		}
	}

	return opts
}

type LanguageMatch struct {
	language       string
	matchedPattern string
	contents       string
}
