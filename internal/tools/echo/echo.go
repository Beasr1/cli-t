// internal/tools/echo/echo.go
package echo

import (
	"context"
	"fmt"
	"strings"

	"cli-t/internal/command"
)

type Command struct {
	// command-specific fields could go here
	noNewline bool
}

func New() command.Command {
	return &Command{}
}

// These methods define how the command appears in CLI-T
func (c *Command) Name() string {
	return "echo"
}

func (c *Command) Usage() string {
	return "[OPTIONS] [STRING...]"
}

func (c *Command) Description() string {
	return "Display a line of text"
}

func (c *Command) ValidateArgs(args []string) error {
	// Echo accepts any arguments
	return nil
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	// Join all positional arguments with spaces
	text := strings.Join(args.Positional, " ")

	// Check if we should add newline
	// TODO: This would come from flags once implemented
	if c.noNewline {
		fmt.Fprint(args.Stdout, text)
	} else {
		fmt.Fprintln(args.Stdout, text)
	}

	return nil
}
