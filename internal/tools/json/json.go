package json

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/file"
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "json"
}

func (c *Command) Usage() string {
	return "[file]"
}

func (c *Command) Description() string {
	return "Validate JSON file or stdin"
}

func (c *Command) ValidateArgs(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("too many arguments, expected 0 or 1 file")
	}
	return nil
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	// Read input
	source, content, err := c.readInput(args)
	if err != nil {
		return err
	}

	logger.Debug("Validating JSON", "source", source, "size", len(content))

	// Validate JSON
	if err := ValidateJSON(content); err != nil {
		// Invalid JSON - write to stderr and return error
		fmt.Fprintf(args.Stderr, "Invalid JSON in %s: %v\n", source, err)
		return fmt.Errorf("validation failed")
	}

	// Valid JSON
	fmt.Fprintf(args.Stdout, "Valid JSON\n")
	return nil
}

// readInput reads from stdin or file, returning source name and content
func (c *Command) readInput(args *command.Args) (source string, content string, err error) {
	if len(args.Positional) == 0 {
		// Read from stdin
		logger.Verbose("Reading from stdin")
		content, err = file.ReadFrom(args.Stdin)
		if err != nil {
			return "", "", fmt.Errorf("error reading from stdin: %w", err)
		}
		return "stdin", content, nil
	}

	// Read from file
	filePath := args.Positional[0]
	logger.Verbose("Reading from file", "path", filePath)

	content, err = file.ReadContent(filePath)
	if err != nil {
		return "", "", fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return filePath, content, nil
}
