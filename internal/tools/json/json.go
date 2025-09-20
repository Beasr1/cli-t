package json

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/file"
	"context"
	"fmt"
	"os"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "json"
}

func (c *Command) Usage() string {
	return "[FILE]"
}

func (c *Command) Description() string {
	return "Validate JSON file or stdin"
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

	// Validate JSON (Inbuilt function)
	// var js json.RawMessage
	// if err := json.Unmarshal([]byte(content), &js); err != nil {
	// 	// Invalid JSON
	// 	fmt.Fprintf(args.Stderr, "Invalid JSON in %s: %v\n", source, err)
	// 	os.Exit(1)
	// }

	if err := ValidateJSON(content); err != nil {
		// Invalid JSON
		fmt.Fprintf(args.Stderr, "Invalid JSON in %s: %v\n", source, err)
		return nil
	}

	// Valid JSON
	fmt.Fprintf(args.Stdout, "Valid JSON\n")
	return nil
}
