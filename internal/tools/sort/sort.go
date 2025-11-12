// internal/tools/sort/sort.go
package sort

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"context"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "sort"
}

func (c *Command) Usage() string {
	return "sort <file>"
}

func (c *Command) Description() string {
	return "sort or merge records (lines) of text and binary files"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "algorithm",
			Shorthand: "a",
			Usage:     "sorting algorithm to use",
			Type:      "string",
			Default:   "",
		},
		{
			Name:      "output",
			Shorthand: "o",
			Usage:     "Output file path (default: stdout or auto-generated)",
			Type:      "string",
			Default:   "",
		},
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	_, output := c.parseFlags(args.Flags)

	inputPath, content, err := io.ReadInput(args)
	if err != nil {
		return err
	}

	// TODO: process
	processed := content

	io.WriteOutput([]byte(processed), inputPath, output, args.Stdout)

	return nil
}

func (c *Command) parseFlags(flags map[string]interface{}) (string, string) {
	algorithm, _ := flags["algorithm"].(string)
	output, _ := flags["output"].(string)

	logger.Debug("Flags processing",
		"algorithm", algorithm,
		"output", output,
	)

	return algorithm, output
}
