package helpers_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"cli-t/internal/command"
	"cli-t/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock command for testing
type mockCommand struct {
	executeFunc func(ctx context.Context, args *command.Args) error
}

func (m *mockCommand) Name() string                     { return "mock" }
func (m *mockCommand) Usage() string                    { return "mock [args]" }
func (m *mockCommand) Description() string              { return "Mock command for testing" }
func (m *mockCommand) ValidateArgs(args []string) error { return nil }
func (m *mockCommand) Execute(ctx context.Context, args *command.Args) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, args)
	}
	return nil
}

func TestHelpers_TestArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		validate func(t *testing.T, args *command.Args)
	}{
		{
			name: "no arguments",
			args: []string{},
			validate: func(t *testing.T, args *command.Args) {
				assert.Empty(t, args.Positional)
				assert.NotNil(t, args.Flags)
				assert.NotNil(t, args.Stdin)
				assert.NotNil(t, args.Stdout)
				assert.NotNil(t, args.Stderr)
				assert.NotNil(t, args.Env)
				assert.NotNil(t, args.Config)
			},
		},
		{
			name: "with arguments",
			args: []string{"file1.txt", "file2.txt", "--flag"},
			validate: func(t *testing.T, args *command.Args) {
				assert.Equal(t, []string{"file1.txt", "file2.txt", "--flag"}, args.Positional)
			},
		},
		{
			name: "config defaults",
			args: []string{},
			validate: func(t *testing.T, args *command.Args) {
				assert.False(t, args.Config.Verbose)
				assert.False(t, args.Config.Debug)
				assert.True(t, args.Config.NoColor)
				assert.Equal(t, "plain", args.Config.Output)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := helpers.TestArgs(tt.args...)
			tt.validate(t, args)
		})
	}
}

func TestHelpers_ExecuteCommand(t *testing.T) {
	tests := []struct {
		name       string
		cmd        command.Command
		args       []string
		wantStdout string
		wantStderr string
		wantErr    bool
	}{
		{
			name: "successful execution with stdout",
			cmd: &mockCommand{
				executeFunc: func(ctx context.Context, args *command.Args) error {
					_, err := args.Stdout.Write([]byte("Hello from stdout"))
					return err
				},
			},
			args:       []string{"arg1", "arg2"},
			wantStdout: "Hello from stdout",
			wantStderr: "",
			wantErr:    false,
		},
		{
			name: "successful execution with stderr",
			cmd: &mockCommand{
				executeFunc: func(ctx context.Context, args *command.Args) error {
					_, err := args.Stderr.Write([]byte("Warning message"))
					return err
				},
			},
			args:       []string{},
			wantStdout: "",
			wantStderr: "Warning message",
			wantErr:    false,
		},
		{
			name: "execution with error",
			cmd: &mockCommand{
				executeFunc: func(ctx context.Context, args *command.Args) error {
					return errors.New("command failed")
				},
			},
			args:       []string{},
			wantStdout: "",
			wantStderr: "",
			wantErr:    true,
		},
		{
			name: "execution with both output and error",
			cmd: &mockCommand{
				executeFunc: func(ctx context.Context, args *command.Args) error {
					args.Stdout.Write([]byte("Partial output"))
					args.Stderr.Write([]byte("Error occurred"))
					return errors.New("failed after output")
				},
			},
			args:       []string{},
			wantStdout: "Partial output",
			wantStderr: "Error occurred",
			wantErr:    true,
		},
		{
			name: "reading from stdin",
			cmd: &mockCommand{
				executeFunc: func(ctx context.Context, args *command.Args) error {
					// Read from stdin and echo to stdout
					buf := make([]byte, 1024)
					n, err := args.Stdin.Read(buf)
					if err != nil && err.Error() != "EOF" {
						return err
					}
					args.Stdout.Write(buf[:n])
					return nil
				},
			},
			args:       []string{},
			wantStdout: "", // stdin is empty by default
			wantStderr: "",
			wantErr:    false,
		},
		{
			name: "accessing positional args",
			cmd: &mockCommand{
				executeFunc: func(ctx context.Context, args *command.Args) error {
					// Echo positional args to stdout
					output := strings.Join(args.Positional, " ")
					args.Stdout.Write([]byte(output))
					return nil
				},
			},
			args:       []string{"hello", "world", "test"},
			wantStdout: "hello world test",
			wantStderr: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := helpers.ExecuteCommand(t, tt.cmd, tt.args...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantStdout, stdout)
			assert.Equal(t, tt.wantStderr, stderr)
		})
	}
}

func TestHelpers_BufferTypes(t *testing.T) {
	// Verify that the I/O streams are actually bytes.Buffer
	args := helpers.TestArgs()

	// Should be able to type assert to *bytes.Buffer
	stdoutBuf, ok := args.Stdout.(*bytes.Buffer)
	assert.True(t, ok, "Stdout should be *bytes.Buffer")
	assert.NotNil(t, stdoutBuf)

	stderrBuf, ok := args.Stderr.(*bytes.Buffer)
	assert.True(t, ok, "Stderr should be *bytes.Buffer")
	assert.NotNil(t, stderrBuf)

	stdinBuf, ok := args.Stdin.(*bytes.Buffer)
	assert.True(t, ok, "Stdin should be *bytes.Buffer")
	assert.NotNil(t, stdinBuf)
}

func TestHelpers_ModifyingTestArgs(t *testing.T) {
	// Test that we can modify TestArgs for specific test scenarios
	args := helpers.TestArgs("initial", "args")

	// Modify stdin
	stdinBuf := args.Stdin.(*bytes.Buffer)
	stdinBuf.WriteString("test input data")

	// Modify flags
	args.Flags["verbose"] = true
	args.Flags["count"] = 42

	// Modify env
	args.Env["TEST_VAR"] = "test_value"
	args.Env["PATH"] = "/test/path"

	// Modify config
	args.Config.Verbose = true
	args.Config.Debug = true
	args.Config.Output = "json"

	// Verify modifications
	assert.Equal(t, []string{"initial", "args"}, args.Positional)
	assert.Equal(t, true, args.Flags["verbose"])
	assert.Equal(t, 42, args.Flags["count"])
	assert.Equal(t, "test_value", args.Env["TEST_VAR"])
	assert.True(t, args.Config.Verbose)
	assert.True(t, args.Config.Debug)
	assert.Equal(t, "json", args.Config.Output)

	// Read from modified stdin
	data, err := io.ReadAll(args.Stdin)
	require.NoError(t, err)
	assert.Equal(t, "test input data", string(data))
}

// Example of using helpers in actual tests
func ExampleTestArgs() {
	// Create test arguments
	args := helpers.TestArgs("file1.txt", "file2.txt")

	// Add stdin data
	args.Stdin.(*bytes.Buffer).WriteString("input data")

	// Set flags
	args.Flags["lines"] = true
	args.Flags["words"] = false

	// Use in command execution
	cmd := &mockCommand{
		executeFunc: func(ctx context.Context, args *command.Args) error {
			// Command implementation
			return nil
		},
	}

	_ = cmd.Execute(context.Background(), args)

	// Check output
	stdout := args.Stdout.(*bytes.Buffer).String()
	stderr := args.Stderr.(*bytes.Buffer).String()

	_ = stdout
	_ = stderr
}

// Benchmark helper creation
func BenchmarkTestArgs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		args := helpers.TestArgs("arg1", "arg2", "arg3")
		_ = args
	}
}

func BenchmarkExecuteCommand(b *testing.B) {
	cmd := &mockCommand{
		executeFunc: func(ctx context.Context, args *command.Args) error {
			args.Stdout.Write([]byte("test output"))
			return nil
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdout, stderr, err := helpers.ExecuteCommand(&testing.T{}, cmd, "test")
		_ = stdout
		_ = stderr
		_ = err
	}
}
