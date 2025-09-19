package wc

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/file"
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
	return "[--byte|-c] [--line|-l] [--word|-w] [--char|-m]"
}

func (c *Command) Description() string {
	return "Show word, line, character, and byte count information"
}

func (c *Command) ValidateArgs(args []string) error {
	// No validation needed - flags are optional
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
	// Ensure at least one positional argument is provided (the file path)
	if len(args.Positional) < 1 {
		return fmt.Errorf("file path is required")
	}

	// First positional argument is treated as the file path
	filePath := args.Positional[0]

	// --- Parse CLI Flags (these may or may not be present) ---

	// Type assert flags to bool; will be 'false' if unset or type mismatch
	// Note: if user does not pass a flag, the value may be nil, hence the safe '_'
	showBytes, _ := args.Flags["byte"].(bool)
	showLines, _ := args.Flags["line"].(bool)
	showWords, _ := args.Flags["word"].(bool)
	showChars, _ := args.Flags["char"].(bool)

	// If no specific flag is passed, show all counts by default (like Unix `wc`)
	if !showBytes && !showLines && !showWords && !showChars {
		showBytes, showLines, showWords, showChars = true, true, true, true
	}

	// --- File Reading ---

	// Read the full content of the file as a string
	content, err := file.ReadContent(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// --- Perform and Display Requested Counts ---

	// Show line count
	if showLines {
		fmt.Printf("Lines: %d\n", len(strings.Split(strings.TrimSuffix(content, "\n"), "\n")))
	}

	// Show word count
	if showWords {
		fmt.Printf("Words: %d\n", len(strings.Fields(content)))
	}

	// Show byte count (number of bytes in the raw file)
	if showBytes {
		fmt.Printf("Bytes: %d\n", len([]byte(content)))
	}

	// Show character count (based on Unicode runes)
	if showChars {
		fmt.Printf("Characters: %d\n", len([]rune(content)))
	}

	// Command completed successfully
	return nil
}
