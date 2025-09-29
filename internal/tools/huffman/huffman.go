package huffman

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/file"
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
	"os"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "huffman"
}

func (c *Command) Usage() string {
	return "[FILE]"
}

func (c *Command) Description() string {
	return "Comppress and decompress a file"
}

func (c *Command) ValidateArgs(args []string) error {
	// Accept 0 or 1 arguments
	if len(args) > 1 {
		return fmt.Errorf("too many arguments, expected 0 or 1 file")
	}
	return nil
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	var content string
	var err error
	var source string

	// Determine input source
	if len(args.Positional) == 0 {
		// Read from stdin
		content, err = file.ReadFrom(args.Stdin)
		if err != nil {
			fmt.Fprintf(args.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		source = "stdin"
	} else {
		// Read from file
		filePath := args.Positional[0]
		content, err = file.ReadContent(filePath)
		if err != nil {
			fmt.Fprintf(args.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		source = filePath
	}

	//Compress or decompress the file depending upon the flag
	logger.Debug("abc", "source", source, "content", content)

	return nil
}
