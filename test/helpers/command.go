package helpers

import (
	"bytes"
	"context"
	"testing"

	"cli-t/internal/command"
)

// TestArgs creates Args for testing
func TestArgs(args ...string) *command.Args {
	return &command.Args{
		Positional: args,
		Flags:      make(map[string]interface{}),
		Stdin:      &bytes.Buffer{},
		Stdout:     &bytes.Buffer{},
		Stderr:     &bytes.Buffer{},
		Env:        make(map[string]string),
		Config: &command.Config{
			Verbose: false,
			Debug:   false,
			NoColor: true,
			Output:  "plain",
		},
	}
}

// ExecuteCommand runs a command and captures output
func ExecuteCommand(t *testing.T, cmd command.Command, args ...string) (string, string, error) {
	testArgs := TestArgs(args...)
	err := cmd.Execute(context.Background(), testArgs)

	stdout := testArgs.Stdout.(*bytes.Buffer).String()
	stderr := testArgs.Stderr.(*bytes.Buffer).String()

	return stdout, stderr, err
}
