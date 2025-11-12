package cut

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"context"
	"fmt"
	"strings"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "cut"
}

func (c *Command) Usage() string {
	return "[-b list | -c list | -f list] [-d delim] [-s] [file ...]"
}

func (c *Command) Description() string {
	return "Filters out the selected portions from each line in a file"
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "bytes",
			Shorthand: "b",
			Usage:     "Select only these bytes",
			Type:      "string",
			Default:   "",
		},
		{
			Name:      "characters",
			Shorthand: "c",
			Usage:     "Select only these characters",
			Type:      "string",
			Default:   "",
		},
		{
			Name:      "fields",
			Shorthand: "f",
			Usage:     "Select only these fields",
			Type:      "string",
			Default:   "",
		},
		{
			Name:      "delimiter",
			Shorthand: "d",
			Usage:     "Use DELIM instead of TAB for field delimiter",
			Type:      "string",
			Default:   "\t",
		},
		{
			Name:      "suppress",
			Shorthand: "s",
			Usage:     "Suppress lines with no delimiter characters",
			Type:      "bool",
			Default:   false,
		},
		// Add -w, -n if you want
	}
}

func (c *Command) ValidateArgs(args []string) error {
	// Can't use -b, -c, and -f together
	// At least one of -b, -c, or -f must be specified
	// validateArgs only has positional arguments
	logger.Debug("args", "args", args)
	return nil
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	// 1. Parse and validate flags
	mode, list, delimiter, suppress, err := c.parseFlags(args.Flags)
	if err != nil {
		return err
	}

	logger.Debug("Cut mode determined",
		"mode", mode,
		"list", list,
		"delimiter", delimiter,
		"suppress", suppress,
	)

	// 2. Read input
	_, content, err := io.ReadInput(args)
	if err != nil {
		return err
	}

	// 3. Process lines
	return c.processLines(string(content), mode, list, delimiter, suppress, args.Stdout)
}

func (c *Command) parseFlags(flags map[string]interface{}) (mode, list, delimiter string, suppress bool, err error) {
	bytesList, _ := flags["bytes"].(string)
	charsList, _ := flags["characters"].(string)
	fieldsList, _ := flags["fields"].(string)
	delimiter, _ = flags["delimiter"].(string)
	suppress, _ = flags["suppress"].(bool)

	// Validate mutual exclusivity
	modesSet := 0
	if bytesList != "" {
		mode = "bytes"
		list = bytesList
		modesSet++
	}
	if charsList != "" {
		mode = "characters"
		list = charsList
		modesSet++
	}
	if fieldsList != "" {
		mode = "fields"
		list = fieldsList
		modesSet++
	}

	if modesSet == 0 {
		err = fmt.Errorf("you must specify one of -b, -c, or -f")
		return
	}
	if modesSet > 1 {
		err = fmt.Errorf("only one of -b, -c, or -f can be specified")
		return
	}

	return
}

func (c *Command) processLines(content, mode, list, delimiter string, suppress bool, stdout io.Writer) error {
	// Parse the list specification
	selections, err := ParseList(list)
	if err != nil {
		return fmt.Errorf("invalid list specification: %w", err)
	}

	logger.Debug("Parsed selections", "count", len(selections), "selections", selections)

	// Split content into lines
	lines := strings.Split(content, "\n")

	// Process each line
	for i, line := range lines {
		// Skip empty last line
		if i == len(lines)-1 && line == "" {
			continue
		}

		var extracted string

		switch mode {
		case "fields":
			// Check if line contains delimiter (for -s flag)
			if suppress && !strings.ContainsRune(line, rune(delimiter[0])) {
				// Skip lines with no delimiter when -s is set
				continue
			}

			extracted = ExtractFields(line, rune(delimiter[0]), selections)

		case "characters":
			extracted = ExtractChars(line, selections)

		case "bytes":
			extracted = ExtractBytes(line, selections)

		default:
			return fmt.Errorf("unknown mode: %s", mode)
		}

		// Write result
		fmt.Fprintln(stdout, extracted)
	}

	return nil
}
