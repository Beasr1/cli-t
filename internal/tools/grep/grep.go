// internal/tools/grep/grep.go
package grep

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
	"os"
)

// TODO: Implement custom regex engine
// Currently using Go's regexp package (easy mode).
// Hard mode requires building:
// - Regex parser (convert string pattern to AST)
// - NFA/DFA construction (finite automata)
// - Backtracking or automata-based matching
// - Support for: ., *, +, ?, [], ^, $, |, ()
// Resources:
// - "Regular Expression Matching Can Be Simple And Fast" by Russ Cox
// - https://swtch.com/~rsc/regexp/
// - Dragon Book chapter on lexical analysis

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "grep"
}

func (c *Command) Usage() string {
	return "grep"
}

func (c *Command) Description() string {
	return "The grep utility searches any given input files, selecting lines that match one or more patterns."
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "recurse",
			Shorthand: "r",
			Type:      "bool",
			Default:   "false",
			Usage:     "recurse a directory tree",
		},
		{
			Name:      "invert-match",
			Shorthand: "v",
			Type:      "bool",
			Default:   "false",
			Usage:     "invert the matching",
		},
		{
			Name:      "ignore-case",
			Shorthand: "i",
			Type:      "bool",
			Default:   "false",
			Usage:     "case insensitive matching",
		},
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	recurse, invert, ignoreCase := c.parseFlags(args.Flags)

	pattern := args.Positional[0]
	if ignoreCase {
		pattern = "(?i)" + pattern
	}

	// I need to read streams of data line by line
	matched := false
	if len(args.Positional) == 1 {
		// Search stdin (still streaming with bufio.Scanner!)

		match, err := SearchStdin(args.Stdin, pattern, invert)
		if err != nil {
			return err
		}
		matched = match || matched

	} else if len(args.Positional) >= 2 {
		// Search file(s)
		files := args.Positional[1:]
		showFilename := len(files) > 1 || recurse
		logger.Info("f", "file", files)
		for _, file := range files {

			info, err := os.Stat(file)
			if err != nil {
				return err
			}

			if info.IsDir() {
				// It's a directory - need to walk it
				if !recurse {
					return fmt.Errorf("%s is a directory (use -r to search recursively)", file)
				}

				match, err := SearchDirectory(file, pattern, invert)
				if err != nil {
					return err
				}
				matched = match || matched
			} else {
				// It's a file - search it
				match, err := SearchFile(file, pattern, showFilename, invert)
				if err != nil {
					return err
				}
				matched = match || matched
			}

		}
	}

	logger.Debug("Grepping", "matched", matched)

	// If no matches found, still return nil (success)
	// Grep executed successfully, just didn't find anything
	if !matched {
		// For now, we'll just exit 0
		// Real grep would exit 1, but that requires architecture changes
	}

	// TODO: Implement proper exit codes
	// - Exit 0: Found matches
	// - Exit 1: No matches found
	// - Exit 2: Error occurred
	// Current implementation: Always exits 0 (requires architecture changes)

	return nil // Exit 0
}

func (c *Command) parseFlags(flags map[string]interface{}) (bool, bool, bool) {
	recurse, _ := flags["recurse"].(bool)
	invert, _ := flags["invert-match"].(bool)
	ignoreCase, _ := flags["ignore-case"].(bool)
	return recurse, invert, ignoreCase
}
