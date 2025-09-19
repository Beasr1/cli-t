package wc_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"cli-t/internal/command"
	"cli-t/internal/tools/wc"
	"cli-t/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWC_Metadata(t *testing.T) {
	cmd := wc.New()

	assert.Equal(t, "wc", cmd.Name())
	assert.Equal(t, "[--byte|-c] [--line|-l] [--word|-w] [--char|-m]", cmd.Usage())
	assert.Equal(t, "Show word, line, character, and byte count information", cmd.Description())
}

func TestWC_ValidateArgs(t *testing.T) {
	cmd := wc.New()

	// wc accepts any arguments
	assert.NoError(t, cmd.ValidateArgs([]string{}))
	assert.NoError(t, cmd.ValidateArgs([]string{"file.txt"}))
	assert.NoError(t, cmd.ValidateArgs([]string{"file1.txt", "file2.txt"}))
}

func TestWC_Execute(t *testing.T) {
	// Create test files
	tmpDir := t.TempDir()

	// Test file with known content
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "Hello world\nThis is a test\nFile with multiple lines\n"
	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Empty file
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	err = os.WriteFile(emptyFile, []byte(""), 0644)
	require.NoError(t, err)

	// Unicode file
	unicodeFile := filepath.Join(tmpDir, "unicode.txt")
	unicodeContent := "Hello 世界\n🌍🌎🌏\n"
	err = os.WriteFile(unicodeFile, []byte(unicodeContent), 0644)
	require.NoError(t, err)

	tests := []struct {
		name       string
		file       string
		flags      map[string]interface{}
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "count all (default)",
			file:       testFile,
			flags:      map[string]interface{}{},
			wantOutput: "Lines: 3\nWords: 9\nBytes: 48\nCharacters: 48\n",
		},
		{
			name: "count lines only",
			file: testFile,
			flags: map[string]interface{}{
				"line": true,
			},
			wantOutput: "Lines: 3\n",
		},
		{
			name: "count words only",
			file: testFile,
			flags: map[string]interface{}{
				"word": true,
			},
			wantOutput: "Words: 9\n",
		},
		{
			name: "count bytes only",
			file: testFile,
			flags: map[string]interface{}{
				"byte": true,
			},
			wantOutput: "Bytes: 48\n",
		},
		{
			name: "count characters",
			file: testFile,
			flags: map[string]interface{}{
				"char": true,
			},
			wantOutput: "Characters: 48\n", // Same as bytes for ASCII
		},
		{
			name: "multiple flags",
			file: testFile,
			flags: map[string]interface{}{
				"line": true,
				"word": true,
			},
			wantOutput: "Lines: 3\nWords: 9\n",
		},
		{
			name: "empty file",
			file: emptyFile,
			flags: map[string]interface{}{
				"line": true,
			},
			wantOutput: "Lines: 1\n", // Empty file has 1 line
		},
		{
			name: "unicode characters",
			file: unicodeFile,
			flags: map[string]interface{}{
				"char": true,
			},
			wantOutput: "Characters: 13\n", // Unicode chars counted correctly
		},
		{
			name: "unicode bytes vs chars",
			file: unicodeFile,
			flags: map[string]interface{}{
				"byte": true,
			},
			wantOutput: "Bytes: 23\n", // More bytes than characters due to UTF-8
		},
		{
			name:    "missing file",
			file:    filepath.Join(tmpDir, "nonexistent.txt"),
			flags:   map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := wc.New()
			args := helpers.TestArgs(tt.file)
			args.Flags = tt.flags

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

func TestWC_NoArgs(t *testing.T) {
	cmd := wc.New()
	args := helpers.TestArgs() // No file argument

	err := cmd.Execute(context.Background(), args)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file path is required")
}

func TestWC_DefineFlags(t *testing.T) {
	cmd := wc.New()

	// Check if command implements FlagDefiner
	flagDefiner, ok := cmd.(command.FlagDefiner)
	require.True(t, ok, "wc command should implement FlagDefiner")

	flags := flagDefiner.DefineFlags()

	// Should have 4 flags
	assert.Len(t, flags, 4)

	// Check each flag
	flagMap := make(map[string]command.Flag)
	for _, f := range flags {
		flagMap[f.Name] = f
	}

	// Check byte flag
	byteFlag, ok := flagMap["byte"]
	assert.True(t, ok)
	assert.Equal(t, "c", byteFlag.Shorthand)
	assert.Equal(t, "bool", byteFlag.Type)
	assert.Equal(t, false, byteFlag.Default)

	// Check line flag
	lineFlag, ok := flagMap["line"]
	assert.True(t, ok)
	assert.Equal(t, "l", lineFlag.Shorthand)
	assert.Equal(t, "bool", lineFlag.Type)

	// Check word flag
	wordFlag, ok := flagMap["word"]
	assert.True(t, ok)
	assert.Equal(t, "w", wordFlag.Shorthand)

	// Check char flag
	charFlag, ok := flagMap["char"]
	assert.True(t, ok)
	assert.Equal(t, "m", charFlag.Shorthand)
}

func TestWC_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		content    string
		flags      map[string]interface{}
		wantOutput string
	}{
		{
			name:       "only newlines",
			content:    "\n\n\n",
			flags:      map[string]interface{}{"line": true},
			wantOutput: "Lines: 3\n",
		},
		{
			name:       "no trailing newline",
			content:    "hello world",
			flags:      map[string]interface{}{"line": true},
			wantOutput: "Lines: 1\n",
		},
		{
			name:       "only spaces",
			content:    "   ",
			flags:      map[string]interface{}{"word": true},
			wantOutput: "Words: 0\n",
		},
		{
			name:       "mixed whitespace",
			content:    "hello\tworld\n  test  \n",
			flags:      map[string]interface{}{"word": true},
			wantOutput: "Words: 3\n",
		},
		{
			name:       "emoji characters",
			content:    "😀😃😄😁",
			flags:      map[string]interface{}{"char": true},
			wantOutput: "Characters: 4\n",
		},
		{
			name:       "emoji bytes",
			content:    "😀😃😄😁",
			flags:      map[string]interface{}{"byte": true},
			wantOutput: "Bytes: 16\n", // Each emoji is 4 bytes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, "test.txt")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			cmd := wc.New()
			args := helpers.TestArgs(testFile)
			args.Flags = tt.flags

			err = cmd.Execute(context.Background(), args)
			assert.NoError(t, err)

			stdout := args.Stdout.(*bytes.Buffer).String()
			assert.Equal(t, tt.wantOutput, stdout)
		})
	}
}

// Benchmark
func BenchmarkWC_LargeFile(b *testing.B) {
	// Create a large test file
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "large.txt")

	// Generate 10MB of text
	var content bytes.Buffer
	for i := 0; i < 100000; i++ {
		content.WriteString("This is a line of text for benchmarking word count performance\n")
	}

	err := os.WriteFile(testFile, content.Bytes(), 0644)
	require.NoError(b, err)

	cmd := wc.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		args := helpers.TestArgs(testFile)
		cmd.Execute(context.Background(), args)
	}
}
