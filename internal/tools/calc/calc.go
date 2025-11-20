// internal/tools/calc/calc.go
package calc

import (
	"cli-t/internal/command"
	"cli-t/internal/shared/io"
	"cli-t/internal/shared/logger"
	"cli-t/internal/tools/calc/stack"
	"cli-t/internal/tools/calc/token"
	"context"
	"fmt"
	"strings"
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

Negative numbers require -- separator (./cli-t calc -- "-5+3")
Floating point precision issues (not rounded)
No function support (sin, cos, sqrt, etc.)
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

	result, err := c.calc(expression)
	if err != nil {
		return err
	}

	fmt.Fprintln(args.Stdout, result)
	return err
}

func (c *Command) parseFlags(flags map[string]interface{}) (string, bool) {
	method, _ := flags["method"].(string)
	fromFile, _ := flags["fromFile"].(bool)
	return method, fromFile
}

func (c *Command) calc(content string) (string, error) {
	tokens, err := token.Tokenizer(content)
	if err != nil {
		return "", err
	}

	logger.Debug("tokenizer", "tokens", tokens)
	postfix, err := stack.ToPostfix(tokens)
	if err != nil {
		return "", err
	}

	logger.Debug("postfix expression", "expression", postfix)
	result, err := stack.EvaluatePostfix(postfix)
	if err != nil {
		return "", err
	}

	// %f: float ,   %g: remove trailing zeros
	logger.Debug("result", "result", result)
	return formatResult(result), nil
}

func formatResult(result float64) string {
	// If result is close to an integer, show as integer
	if result == float64(int64(result)) {
		return fmt.Sprintf("%d", int64(result))
	}
	// Otherwise show with reasonable precision
	formatted := fmt.Sprintf("%.10f", result)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")
	return formatted
}
