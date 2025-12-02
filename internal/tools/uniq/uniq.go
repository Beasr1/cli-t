package uniq

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"context"
	"fmt"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "uniq"
}

func (c *Command) Usage() string {
	return "uniq"
}

func (c *Command) Description() string {
	return "The uniq utility writes a copy of each unique input line"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "count",
			Shorthand: "c",
			Type:      "bool",
			Default:   "false",
			Usage:     "count of the number of times a line appears in the input file",
		},
		{
			Name:      "repeated",
			Shorthand: "d",
			Type:      "bool",
			Default:   "false",
			Usage:     "output only repeated lines",
		},
		{
			Name:      "unique",
			Shorthand: "u",
			Type:      "bool",
			Default:   "false",
			Usage:     "output only uniq lines",
		},
	}
}

// this also reads line by line aaaahahh
// we DON'T need everything IN MEMORY at once
func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	count, repeated, unique := c.parseFlags(args.Flags)

	// Validate flags
	if repeated && unique {
		return fmt.Errorf("options -d and -u are mutually exclusive")
	}

	// Determine input path from args
	var inputPath string
	if len(args.Positional) >= 1 && args.Positional[0] != "-" {
		inputPath = args.Positional[0]
	}

	// Determine output path from args
	var outputPath string
	if len(args.Positional) >= 2 {
		outputPath = args.Positional[1]
	}

	// Setup INPUT reader
	reader, cleanupReader, err := io.GetInputReader(args.Stdin, inputPath)
	if err != nil {
		return err
	}
	defer cleanupReader()

	// Setup OUTPUT writer
	writer, cleanupWriter, err := io.GetOutputWriter(args.Stdout, outputPath)
	if err != nil {
		return err
	}
	defer cleanupWriter()

	opts := Options{
		Count:    count,
		Repeated: repeated,
		Unique:   unique,
	}

	// Call Process - it doesn't know what reader/writer are!
	return Process(reader, writer, opts)
}

func (c *Command) parseFlags(flags map[string]interface{}) (bool, bool, bool) {
	count, _ := flags["count"].(bool)
	repeated, _ := flags["repeated"].(bool)
	unique, _ := flags["unique"].(bool)
	return count, repeated, unique
}
