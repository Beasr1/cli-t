// internal/tools/calc/calc.go
package calc

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"cli-t/internal/tools/calc/token"
	"context"
	"fmt"
)

type Command struct{}

func New() command.Command {
	return &Command{}
}

func (c *Command) Name() string {
	return "calc"
}

func (c *Command) Usage() string {
	return "calc []"
}

func (c *Command) Description() string {
	return "calculator"
}

func (c *Command) ValidateArgs(args []string) error {
	return nil
}

func (c *Command) DefineFlags() []command.Flag {
	return []command.Flag{
		{
			Name:      "method",
			Shorthand: "m",
			Type:      "string",
			Default:   "stack",
			Usage:     "Evaluation method: 'stack' or 'tree'",
		},
		{
			Name:      "from-file",
			Shorthand: "f",
			Type:      "bool",
			Default:   false,
			Usage:     "Read expression from file/stdin instead of direct arg",
		},
	}
}

/*
Cool things
go run ./cmd/cli-t/main.go calc --  "-2+-3"
*/
func (c *Command) Execute(ctx context.Context, args *command.Args) error {
	_, fromFile := c.parseFlags(args.Flags)

	var expression string
	var err error

	if fromFile {
		// Use existing io.ReadInput for file/stdin
		_, content, err := io.ReadInput(args)
		if err != nil {
			return err
		}
		expression = string(content)
	} else {
		// Default: positional[0] is the expression
		if len(args.Positional) == 0 {
			return fmt.Errorf("no expression provided")
		}
		expression = args.Positional[0]
	}

	// Now evaluate based on method
	_ = c.calc(expression)

	// err = io.WriteOutput(result, inputPath, "", args.Stdout, c.Name())
	return err
}

func (c *Command) parseFlags(flags map[string]interface{}) (string, bool) {
	method, _ := flags["method"].(string)
	fromFile, _ := flags["fromFile"].(bool)
	return method, fromFile
}

func (c *Command) calc(content string) float64 {
	var result float64
	// if method == "tree" {
	// 	result, err = c.evaluateWithTree(expression)
	// } else {
	// 	result, err = c.evaluateWithStack(expression)
	// }

	tokens, err := token.Tokenizer(content)
	if err != nil {
		logger.Error("error", "err", err)
	}
	logger.Info("tokens", "tokens", tokens)
	return result
}
