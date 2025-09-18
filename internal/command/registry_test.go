package command_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"

	"cli-t/internal/command"
	"cli-t/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock command for testing
type mockCommand struct {
	name         string
	usage        string
	description  string
	executeFunc  func(ctx context.Context, args *command.Args) error
	validateFunc func(args []string) error
}

func (m *mockCommand) Name() string        { return m.name }
func (m *mockCommand) Usage() string       { return m.usage }
func (m *mockCommand) Description() string { return m.description }

func (m *mockCommand) Execute(ctx context.Context, args *command.Args) error {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, args)
	}
	return nil
}

func (m *mockCommand) ValidateArgs(args []string) error {
	if m.validateFunc != nil {
		return m.validateFunc(args)
	}
	return nil
}

func TestRegistry_Register(t *testing.T) {
	// Clear registry before tests (assuming we add a Clear method)
	// For now, we'll test with the global registry

	tests := []struct {
		name    string
		cmd     command.Command
		wantErr bool
		errMsg  string
	}{
		{
			name: "register new command",
			cmd: &mockCommand{
				name:        "test-cmd",
				usage:       "test-cmd [options]",
				description: "A test command",
			},
			wantErr: false,
		},
		{
			name: "register duplicate command",
			cmd: &mockCommand{
				name:        "test-cmd", // Same name as above
				usage:       "test-cmd [options]",
				description: "Another test command",
			},
			wantErr: true,
			errMsg:  "command test-cmd already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := command.Register(tt.cmd)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegistry_Get(t *testing.T) {
	// Register a test command
	testCmd := &mockCommand{
		name:        "get-test",
		usage:       "get-test",
		description: "Test for Get",
	}
	require.NoError(t, command.Register(testCmd))

	tests := []struct {
		name     string
		cmdName  string
		wantCmd  bool
		validate func(t *testing.T, cmd command.Command)
	}{
		{
			name:    "get existing command",
			cmdName: "get-test",
			wantCmd: true,
			validate: func(t *testing.T, cmd command.Command) {
				assert.Equal(t, "get-test", cmd.Name())
				assert.Equal(t, "Test for Get", cmd.Description())
			},
		},
		{
			name:    "get non-existing command",
			cmdName: "non-existent",
			wantCmd: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, ok := command.Get(tt.cmdName)
			assert.Equal(t, tt.wantCmd, ok)
			if tt.wantCmd && tt.validate != nil {
				tt.validate(t, cmd)
			}
		})
	}
}

func TestRegistry_List(t *testing.T) {
	// Note: This test might be affected by previous tests
	// In a real scenario, we'd want to isolate the registry

	// Register some commands
	cmds := []command.Command{
		&mockCommand{name: "alpha", usage: "alpha", description: "Alpha command"},
		&mockCommand{name: "beta", usage: "beta", description: "Beta command"},
		&mockCommand{name: "gamma", usage: "gamma", description: "Gamma command"},
	}

	for _, cmd := range cmds {
		_ = command.Register(cmd) // Ignore errors for existing commands
	}

	list := command.List()

	// Should be sorted alphabetically
	assert.Contains(t, list, "alpha")
	assert.Contains(t, list, "beta")
	assert.Contains(t, list, "gamma")

	// Verify sorting
	for i := 1; i < len(list); i++ {
		assert.True(t, list[i-1] < list[i], "List should be sorted")
	}
}

func TestRegistry_GetAll(t *testing.T) {
	// Register a command
	testCmd := &mockCommand{
		name:        "getall-test",
		usage:       "getall-test",
		description: "Test for GetAll",
	}
	require.NoError(t, command.Register(testCmd))

	all := command.GetAll()

	// Should contain our test command
	cmd, exists := all["getall-test"]
	assert.True(t, exists)
	assert.Equal(t, "getall-test", cmd.Name())

	// Verify it's a copy (modifying the returned map shouldn't affect registry)
	delete(all, "getall-test")

	// Original should still exist
	_, exists = command.Get("getall-test")
	assert.True(t, exists)
}

func TestCommand_Execute(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *mockCommand
		args    *command.Args
		wantErr bool
		check   func(t *testing.T, args *command.Args)
	}{
		{
			name: "successful execution",
			cmd: &mockCommand{
				name: "exec-test",
				executeFunc: func(ctx context.Context, args *command.Args) error {
					// Write to stdout
					_, err := args.Stdout.Write([]byte("Hello, World!"))
					return err
				},
			},
			args:    helpers.TestArgs(),
			wantErr: false,
			check: func(t *testing.T, args *command.Args) {
				stdout := args.Stdout.(*bytes.Buffer).String()
				assert.Equal(t, "Hello, World!", stdout)
			},
		},
		{
			name: "execution with error",
			cmd: &mockCommand{
				name: "exec-error-test",
				executeFunc: func(ctx context.Context, args *command.Args) error {
					return errors.New("execution failed")
				},
			},
			args:    helpers.TestArgs(),
			wantErr: true,
		},
		{
			name: "execution with context cancellation",
			cmd: &mockCommand{
				name: "exec-cancel-test",
				executeFunc: func(ctx context.Context, args *command.Args) error {
					select {
					case <-ctx.Done():
						return ctx.Err()
					default:
						return nil
					}
				},
			},
			args:    helpers.TestArgs(),
			wantErr: true,
			check: func(t *testing.T, args *command.Args) {
				// Context should be cancelled
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.name == "execution with context cancellation" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel() // Cancel immediately
			}

			err := tt.cmd.Execute(ctx, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.check != nil {
				tt.check(t, tt.args)
			}
		})
	}
}

func TestCommand_ValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		cmd     *mockCommand
		args    []string
		wantErr bool
	}{
		{
			name: "valid args",
			cmd: &mockCommand{
				name: "validate-test",
				validateFunc: func(args []string) error {
					if len(args) == 0 {
						return errors.New("no arguments provided")
					}
					return nil
				},
			},
			args:    []string{"file.txt"},
			wantErr: false,
		},
		{
			name: "invalid args",
			cmd: &mockCommand{
				name: "validate-test",
				validateFunc: func(args []string) error {
					if len(args) == 0 {
						return errors.New("no arguments provided")
					}
					return nil
				},
			},
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.ValidateArgs(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkRegistry_Get(b *testing.B) {
	// Register many commands
	for i := 0; i < 100; i++ {
		cmd := &mockCommand{
			name:        fmt.Sprintf("bench-cmd-%d", i),
			usage:       "bench",
			description: "Benchmark command",
		}
		_ = command.Register(cmd)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = command.Get("bench-cmd-50")
	}
}

func BenchmarkRegistry_List(b *testing.B) {
	// Register many commands
	for i := 0; i < 100; i++ {
		cmd := &mockCommand{
			name:        fmt.Sprintf("bench-list-%d", i),
			usage:       "bench",
			description: "Benchmark command",
		}
		_ = command.Register(cmd)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = command.List()
	}
}
