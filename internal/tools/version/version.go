// internal/tools/version/version.go
package versioncmd

import (
	"context"
	"encoding/json"
	"fmt"

	"cli-t/internal/command"
	"cli-t/pkg/version"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "version"
}

func (c *Command) Usage() string {
	return "[--short|-s] [--json]"
}

func (c *Command) Description() string {
	return "Show version information"
}

func (c *Command) ValidateArgs(args []string) error {
	// No validation needed - flags are optional
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "short",
			Shorthand: "s",
			Usage:     "Show version number only",
			Type:      "bool",
			Default:   false,
		},
		{
			Name:    "json",
			Usage:   "Output version info as JSON",
			Type:    "bool",
			Default: false,
		},
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	// Get flags from parsed args
	short, _ := args.Flags["short"].(bool)
	jsonOut, _ := args.Flags["json"].(bool)

	switch {
	case short:
		fmt.Fprintln(args.Stdout, version.Version)
	case jsonOut:
		info := version.GetInfo()
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal version info: %w", err)
		}
		fmt.Fprintln(args.Stdout, string(data))
	default:
		fmt.Fprintln(args.Stdout, version.DetailedString())
	}

	return nil
}
