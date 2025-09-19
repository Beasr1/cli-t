package echo_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"cli-t/internal/tools/echo"
	"cli-t/test/helpers"

	"github.com/stretchr/testify/assert"
)

func TestEcho_Metadata(t *testing.T) {
	cmd := echo.New()

	assert.Equal(t, "echo", cmd.Name())
	assert.Equal(t, "[OPTIONS] [STRING...]", cmd.Usage())
	assert.Equal(t, "Display a line of text", cmd.Description())
}

func TestEcho_ValidateArgs(t *testing.T) {
	cmd := echo.New()

	// Echo accepts any arguments
	assert.NoError(t, cmd.ValidateArgs([]string{}))
	assert.NoError(t, cmd.ValidateArgs([]string{"hello"}))
	assert.NoError(t, cmd.ValidateArgs([]string{"hello", "world", "test"}))
}

func TestEcho_Execute(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "single word",
			args:       []string{"hello"},
			wantOutput: "hello\n",
		},
		{
			name:       "multiple words",
			args:       []string{"hello", "world"},
			wantOutput: "hello world\n",
		},
		{
			name:       "empty args",
			args:       []string{},
			wantOutput: "\n",
		},
		{
			name:       "special characters",
			args:       []string{"hello\tworld", "test\nline"},
			wantOutput: "hello\tworld test\nline\n",
		},
		{
			name:       "numbers and symbols",
			args:       []string{"123", "!@#", "$%^"},
			wantOutput: "123 !@# $%^\n",
		},
		{
			name:       "unicode",
			args:       []string{"Hello", "世界", "🌍"},
			wantOutput: "Hello 世界 🌍\n",
		},
		{
			name:       "quoted strings",
			args:       []string{"\"hello world\"", "'test'"},
			wantOutput: "\"hello world\" 'test'\n",
		},
		{
			name:       "many arguments",
			args:       []string{"one", "two", "three", "four", "five"},
			wantOutput: "one two three four five\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := echo.New()
			args := helpers.TestArgs(tt.args...)

			err := cmd.Execute(context.Background(), args)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				stdout := args.Stdout.(*bytes.Buffer).String()
				assert.Equal(t, tt.wantOutput, stdout)
			}
		})
	}
}

func TestEcho_LongInput(t *testing.T) {
	// Test with very long input
	var longArgs []string
	for i := 0; i < 1000; i++ {
		longArgs = append(longArgs, "word")
	}

	cmd := echo.New()
	args := helpers.TestArgs(longArgs...)

	err := cmd.Execute(context.Background(), args)
	assert.NoError(t, err)

	stdout := args.Stdout.(*bytes.Buffer).String()
	// Should have 999 spaces between 1000 words, plus newline
	assert.Contains(t, stdout, "word word")
	assert.True(t, strings.HasSuffix(stdout, "\n"))
	assert.Equal(t, 999, strings.Count(stdout, " "))
}

// TODO: Add tests for -n flag when implemented
func TestEcho_NoNewlineFlag(t *testing.T) {
	t.Skip("Waiting for flag implementation")

	// When flags are implemented:
	// cmd := echo.New()
	// args := helpers.TestArgs("hello")
	// args.Flags["n"] = true
	//
	// err := cmd.Execute(context.Background(), args)
	// assert.NoError(t, err)
	//
	// stdout := args.Stdout.(*bytes.Buffer).String()
	// assert.Equal(t, "hello", stdout) // No newline
}

// Benchmark
func BenchmarkEcho_Simple(b *testing.B) {
	cmd := echo.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := helpers.TestArgs("hello", "world")
		cmd.Execute(context.Background(), args)
	}
}

func BenchmarkEcho_ManyArgs(b *testing.B) {
	cmd := echo.New()
	testArgs := make([]string, 100)
	for i := range testArgs {
		testArgs[i] = "test"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := helpers.TestArgs(testArgs...)
		cmd.Execute(context.Background(), args)
	}
}
