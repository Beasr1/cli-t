// internal/tools/config/config.go
package config

import (
	"context"
	"fmt"

	"cli-t/internal/command"
	"cli-t/internal/config"

	"gopkg.in/yaml.v3"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "config"
}

func (c *Command) Usage() string {
	return "[show|edit]"
}

func (c *Command) Description() string {
	return "Manage CLI-T configuration"
}

func (c *Command) ValidateArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand required: show or edit")
	}

	switch args[0] {
	case "show", "edit":
		return nil
	default:
		return fmt.Errorf("unknown subcommand: %s", args[0])
	}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	if len(args.Positional) == 0 {
		return fmt.Errorf("subcommand required: show or edit")
	}

	switch args.Positional[0] {
	case "show":
		return c.show(args)
	case "edit":
		return c.edit(args)
	default:
		return fmt.Errorf("unknown subcommand: %s", args.Positional[0])
	}
}

func (c *Command) show(args *command.Args) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Pretty print as YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	fmt.Fprint(args.Stdout, string(data))
	return nil
}

func (c *Command) edit(args *command.Args) error {
	// TODO: Open config in editor
	return fmt.Errorf("not implemented yet")
}
