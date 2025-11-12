package wc

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
	"strings"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "wc"
}

func (c *Command) Usage() string {
	return "[--byte|-c] [--line|-l] [--word|-w] [--char|-m] [file]"
}

func (c *Command) Description() string {
	return "Show word, line, character, and byte count information"
}

func (c *Command) ValidateArgs(args []string) error {
	// Accept 0 or 1 file
	if len(args) > 1 {
		return fmt.Errorf("too many arguments, expected 0 or 1 file")
	}
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "byte",
			Shorthand: "c",
			Usage:     "Show number of Bytes only",
			Type:      "bool",
			Default:   false,
		},
		{
			Name:      "line",
			Shorthand: "l",
			Usage:     "Show number of lines only",
			Type:      "bool",
			Default:   false,
		},
		{
			Name:      "word",
			Shorthand: "w",
			Usage:     "Show number of words only",
			Type:      "bool",
			Default:   false,
		},
		{
			Name:      "char",
			Shorthand: "m",
			Usage:     "Show number of characters only",
			Type:      "bool",
			Default:   false,
		},
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	// Parse flags
	showBytes, showLines, showWords, showChars := c.parseFlags(args.Flags)

	logger.Debug("WC configuration",
		"bytes", showBytes,
		"lines", showLines,
		"words", showWords,
		"chars", showChars,
	)

	// Read input
	_, content, err := io.ReadInput(args)
	if err != nil {
		return err
	}

	// Display counts
	return c.displayCounts(string(content), showBytes, showLines, showWords, showChars, args.Stdout)
}

// parseFlags extracts flag values
func (c *Command) parseFlags(flags map[string]interface{}) (showBytes, showLines, showWords, showChars bool) {
	// --- Parse CLI Flags (these may or may not be present) ---

	// Type assert flags to bool; will be 'false' if unset or type mismatch
	// Note: if user does not pass a flag, the value may be nil, hence the safe '_'
	showBytes, _ = flags["byte"].(bool)
	showLines, _ = flags["line"].(bool)
	showWords, _ = flags["word"].(bool)
	showChars, _ = flags["char"].(bool)

	// If no specific flag is passed, show all counts by default
	if !showBytes && !showLines && !showWords && !showChars {
		showBytes, showLines, showWords, showChars = true, true, true, true
	}

	return
}

// displayCounts performs counting and outputs results
func (c *Command) displayCounts(content string, showBytes, showLines, showWords, showChars bool, stdout io.Writer) error {
	// Show line count
	if showLines {
		lineCount := len(strings.Split(strings.TrimSuffix(content, "\n"), "\n"))
		fmt.Fprintf(stdout, "Lines: %d\n", lineCount)
	}

	// Show word count
	if showWords {
		wordCount := len(strings.Fields(content))
		fmt.Fprintf(stdout, "Words: %d\n", wordCount)
	}

	// Show byte count
	if showBytes {
		byteCount := len([]byte(content))
		fmt.Fprintf(stdout, "Bytes: %d\n", byteCount)
	}

	// Show character count
	if showChars {
		charCount := len([]rune(content))
		fmt.Fprintf(stdout, "Characters: %d\n", charCount)
	}

	return nil
}
