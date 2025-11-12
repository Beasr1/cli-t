package json

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
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
	source, content, err := io.ReadInput(args)
	if err != nil {
		return err
	}

	logger.Debug("Validating JSON", "source", source, "size", len(content))

	// Validate JSON
	if err := ValidateJSON(string(content)); err != nil {
		// Invalid JSON - write to stderr and return error
		fmt.Fprintf(args.Stderr, "Invalid JSON in %s: %v\n", source, err)
		return fmt.Errorf("validation failed")
	}

	// Valid JSON
	fmt.Fprintf(args.Stdout, "Valid JSON\n")
	return nil
}
