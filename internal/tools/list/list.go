// internal/tools/list/list.go
package list

import (
	"context"
	"fmt"
	"strings"

	"cli-t/internal/command"
	"cli-t/internal/shared/logger"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "list"
}

func (c *Command) Usage() string {
	return ""
}

func (c *Command) Description() string {
	return "List available tools"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	tools := command.List()

	logger.Info("Available tools", "count", len(tools))

	fmt.Println("\nAvailable tools:")
	fmt.Println(strings.Repeat("-", 50))

	for _, name := range tools {
		if cmd, ok := command.Get(name); ok {
			fmt.Fprintf(args.Stdout, "  %-15s %s\n", name, cmd.Description())
		}
	}

	fmt.Printf("\nTotal: %d tools\n", len(tools))
	fmt.Println("\nUse 'cli-t <tool> --help' for tool-specific help")
	return nil
}
