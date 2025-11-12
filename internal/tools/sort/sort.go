// internal/tools/sort/sort.go
package sort

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"context"
	"sort"
	"strings"
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
	algorithm, output := c.parseFlags(args.Flags)

	inputPath, content, err := io.ReadInput(args)
	if err != nil {
		return err
	}

	processed := c.sort(string(content), algorithm)

	io.WriteOutput([]byte(processed), inputPath, output, args.Stdout, c.Name())
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

func (c *Command) sort(content, algorithm string) string {
	// 1. Split into lines
	lines := strings.Split(content, "\n")

	// 2. Handle empty last line
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// 3. Sort the lines
	switch algorithm {
	default:
		logger.Debug("default algorithm")
		sort.Strings(lines)
	}

	// 4. Join back
	return strings.Join(lines, "\n")
}
