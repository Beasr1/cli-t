package urlsh

import (
	"cli-t/internal/command"
	"context"
)

/*
	sankat.
	should I do it for my local to local
	or make it sharable. so need to save the shortened url and actual
	or make a determinitic shortner and lengthening the url with a key of sorts.
	both of us should use same key
*/

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "urlsh"
}

func (c *Command) Usage() string {
	return "urlsh --port <port>"
}

func (c *Command) Description() string {
	return "Shorten URL link"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{}
}

func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	return nil
}
