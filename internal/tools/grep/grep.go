// internal/tools/grep/grep.go
package grep

import (
	"bufio"
	"cli-t/internal/command"
	"cli-t/internal/shared/logger"
	"context"
)

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
	return []command.Flag{}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	_ = c.parseFlags(args.Flags)

	matcher := args.Positional[0]

	// I need to read streams of data line by line
	// In your Execute() method:
	if len(args.Positional) == 0 {
		// Search stdin (still streaming with bufio.Scanner!)
		scanner := bufio.NewScanner(args.Stdin)
	} else {
		// Search file(s)
		files := args.Positional[1:]
		for _, file := range files {
			SearchFile(file, matcher)
		}
	}

	logger.Debug("Grepping", "source", source, "size", len(content))

	return nil
}

func (c *Command) parseFlags(flags map[string]interface{}) int {
	return 0
}
