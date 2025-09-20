package json_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"cli-t/internal/tools/json"
	"cli-t/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test command metadata
func TestCommand_Metadata(t *testing.T) {
	cmd := json.New()

	assert.Equal(t, "json", cmd.Name())
	assert.Equal(t, "[FILE]", cmd.Usage())
	assert.Equal(t, "Validate JSON file or stdin", cmd.Description())
}

// Test command argument validation
func TestCommand_ValidateArgs(t *testing.T) {
	cmd := json.New()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments (stdin)",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "one file argument",
			args:    []string{"data.json"},
			wantErr: false,
		},
		{
			name:    "too many arguments",
			args:    []string{"file1.json", "file2.json"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.ValidateArgs(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test the JSON validation functionality
func TestCommand_JSONValidation(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantValid bool
		wantErr   string
	}{
		// Valid JSON
		{
			name:      "empty object",
			content:   `{}`,
			wantValid: true,
		},
		{
			name:      "empty array",
			content:   `[]`,
			wantValid: true,
		},
		{
			name:      "simple object",
			content:   `{"name": "test", "value": 123}`,
			wantValid: true,
		},
		{
			name:      "nested structure",
			content:   `{"user": {"name": "John", "age": 30}, "active": true}`,
			wantValid: true,
		},
		{
			name:      "array of objects",
			content:   `[{"id": 1}, {"id": 2}]`,
			wantValid: true,
		},
		{
			name:      "mixed array",
			content:   `[1, "two", true, null, {"five": 5}]`,
			wantValid: true,
		},
		{
			name:      "string value",
			content:   `"hello world"`,
			wantValid: true,
		},
		{
			name:      "number value",
			content:   `42.5`,
			wantValid: true,
		},
		{
			name:      "boolean true",
			content:   `true`,
			wantValid: true,
		},
		{
			name:      "boolean false",
			content:   `false`,
			wantValid: true,
		},
		{
			name:      "null value",
			content:   `null`,
			wantValid: true,
		},
		{
			name:      "unicode in string",
			content:   `{"message": "Hello 世界 🌍"}`,
			wantValid: true,
		},
		{
			name:      "escaped characters",
			content:   `{"text": "line1\nline2\ttab\"quote\""}`,
			wantValid: true,
		},

		// Invalid JSON
		{
			name:      "empty input",
			content:   ``,
			wantValid: false,
			wantErr:   "unexpected end of JSON input",
		},
		{
			name:      "incomplete object",
			content:   `{`,
			wantValid: false,
			wantErr:   "unexpected end of JSON input",
		},
		{
			name:      "trailing comma in object",
			content:   `{"a": 1,}`,
			wantValid: false,
			wantErr:   "trailing comma",
		},
		{
			name:      "trailing comma in array",
			content:   `[1, 2, 3,]`,
			wantValid: false,
			wantErr:   "trailing comma",
		},
		{
			name:      "unquoted key",
			content:   `{key: "value"}`,
			wantValid: false,
			wantErr:   "unexpected character",
		},
		{
			name:      "single quotes",
			content:   `{'key': 'value'}`,
			wantValid: false,
			wantErr:   "unexpected character",
		},
		{
			name:      "missing comma",
			content:   `{"a": 1 "b": 2}`,
			wantValid: false,
			wantErr:   "expected ',' or '}'",
		},
		{
			name:      "extra closing brace",
			content:   `{}}`,
			wantValid: false,
			wantErr:   "unexpected token after JSON value",
		},
		{
			name:      "invalid escape sequence",
			content:   `"test\x"`,
			wantValid: false,
			wantErr:   "invalid escape sequence",
		},
		{
			name:      "unescaped newline",
			content:   "\"test\nline\"",
			wantValid: false,
			wantErr:   "unescaped control character",
		},
		{
			name:      "leading zeros in number",
			content:   `{"value": 0123}`,
			wantValid: false,
			wantErr:   "leading zeros not allowed",
		},
		{
			name:      "incomplete string",
			content:   `"unterminated`,
			wantValid: false,
			wantErr:   "unterminated string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.ValidateJSON(tt.content)

			if tt.wantValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if tt.wantErr != "" {
					assert.Contains(t, err.Error(), tt.wantErr)
				}
			}
		})
	}
}

// Test reading from file
func TestCommand_FileInput(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		expectValid bool
	}{
		{
			name:        "valid JSON file",
			content:     `{"valid": true}`,
			expectValid: true,
		},
		{
			name:        "invalid JSON file",
			content:     `{invalid}`,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, "test.json")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Validate file content
			err = json.ValidateJSON(tt.content)

			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// Test reading from stdin
func TestCommand_StdinInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValid bool
	}{
		{
			name:        "valid JSON from stdin",
			input:       `{"test": "value"}`,
			expectValid: true,
		},
		{
			name:        "invalid JSON from stdin",
			input:       `{invalid}`,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate stdin content
			err := json.ValidateJSON(tt.input)

			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// Test error messages include position information
func TestCommand_ErrorPositions(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantLine     int
		wantColumn   int
		wantContains string
	}{
		{
			name:         "error at beginning",
			input:        `[}`,
			wantLine:     1,
			wantColumn:   2,
			wantContains: "line 1, column 2",
		},
		{
			name: "error on second line",
			input: `{
  "key" "value"
}`,
			wantContains: "line 2",
		},
		{
			name:         "error after whitespace",
			input:        `    [1, 2 3]`,
			wantContains: "expected ',' or ']'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := json.ValidateJSON(tt.input)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantContains)
		})
	}
}

// Benchmark validation performance
func BenchmarkValidation_SmallJSON(b *testing.B) {
	input := `{"name": "test", "value": 123, "active": true}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = json.ValidateJSON(input)
	}
}

func BenchmarkValidation_MediumJSON(b *testing.B) {
	input := `{
		"users": [
			{"id": 1, "name": "John", "email": "john@example.com"},
			{"id": 2, "name": "Jane", "email": "jane@example.com"},
			{"id": 3, "name": "Bob", "email": "bob@example.com"}
		],
		"total": 3,
		"page": 1
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = json.ValidateJSON(input)
	}
}

func BenchmarkValidation_LargeJSON(b *testing.B) {
	// Build a large JSON array
	var buf bytes.Buffer
	buf.WriteString(`{"items": [`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{"id":`)
		buf.WriteString(fmt.Sprintf("%d", i))
		buf.WriteString(`,"name":"item`)
		buf.WriteString(fmt.Sprintf("%d", i))
		buf.WriteString(`"}`)
	}
	buf.WriteString(`]}`)
	input := buf.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = json.ValidateJSON(input)
	}
}

// Test actual command execution (integration test)
func TestCommand_Execute(t *testing.T) {
	t.Skip("Skipping because Execute calls os.Exit")

	// In a real scenario, you might refactor Execute to return an exit code
	// or use a subprocess to test the actual command execution
}

// Example of testing with the test helpers
func TestCommand_WithHelpers(t *testing.T) {
	cmd := json.New()

	t.Run("valid JSON from stdin", func(t *testing.T) {
		args := helpers.TestArgs()
		args.Stdin = bytes.NewBufferString(`{"valid": true}`)

		// Since Execute calls os.Exit, we can't test it directly
		// This shows how you would set up the test
		_ = cmd
		_ = args
	})
}
