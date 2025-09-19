package list_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"cli-t/internal/command"
	"cli-t/internal/tools/list"
	"cli-t/test/helpers"

	"github.com/stretchr/testify/assert"
)

func TestList_Metadata(t *testing.T) {
	cmd := list.New()

	assert.Equal(t, "list", cmd.Name())
	assert.Equal(t, "", cmd.Usage())
	assert.Equal(t, "List available tools", cmd.Description())
}

func TestList_ValidateArgs(t *testing.T) {
	cmd := list.New()

	// List accepts any arguments
	assert.NoError(t, cmd.ValidateArgs([]string{}))
	assert.NoError(t, cmd.ValidateArgs([]string{"extra"}))
}

func TestList_Execute(t *testing.T) {
	cmd := list.New()
	args := helpers.TestArgs()

	err := cmd.Execute(context.Background(), args)
	assert.NoError(t, err)

	stdout := args.Stdout.(*bytes.Buffer).String()

	// Check output structure
	assert.Contains(t, stdout, "Available tools:")
	assert.Contains(t, stdout, "----------")
	assert.Contains(t, stdout, "Total:")
	assert.Contains(t, stdout, "Use 'cli-t <tool> --help' for tool-specific help")

	// Check that known tools are listed
	registeredTools := command.List()
	for _, toolName := range registeredTools {
		assert.Contains(t, stdout, toolName)
	}

	// Check format (name and description)
	lines := strings.Split(stdout, "\n")
	foundTool := false
	for _, line := range lines {
		if strings.Contains(line, "echo") && strings.Contains(line, "Display a line of text") {
			foundTool = true
			break
		}
	}
	assert.True(t, foundTool, "Should find echo tool with description")
}

func TestList_ToolCount(t *testing.T) {
	cmd := list.New()
	args := helpers.TestArgs()

	err := cmd.Execute(context.Background(), args)
	assert.NoError(t, err)

	stdout := args.Stdout.(*bytes.Buffer).String()

	// Get actual tool count
	toolCount := len(command.List())

	// Check that count is displayed correctly
	assert.Contains(t, stdout, fmt.Sprintf("Total: %d tools", toolCount))
}

func TestList_Formatting(t *testing.T) {
	cmd := list.New()
	args := helpers.TestArgs()

	err := cmd.Execute(context.Background(), args)
	assert.NoError(t, err)

	stdout := args.Stdout.(*bytes.Buffer).String()
	lines := strings.Split(stdout, "\n")

	// Find tool listing lines (between header and total)
	inTools := false
	var toolLines []string
	for _, line := range lines {
		if strings.Contains(line, "--------") {
			inTools = true
			continue
		}
		if strings.Contains(line, "Total:") {
			break
		}
		if inTools && strings.TrimSpace(line) != "" {
			toolLines = append(toolLines, line)
		}
	}

	// Check that tool lines are properly formatted
	for _, line := range toolLines {
		// Should have at least 2 spaces of indentation
		assert.True(t, strings.HasPrefix(line, "  "))
		// Should have tool name and description separated
		parts := strings.Fields(line)
		assert.True(t, len(parts) >= 2, "Tool line should have name and description")
	}
}

// Benchmark
func BenchmarkList_Execute(b *testing.B) {
	cmd := list.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := helpers.TestArgs()
		cmd.Execute(context.Background(), args)
	}
}
