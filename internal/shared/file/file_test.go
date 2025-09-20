package file_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cli-t/internal/shared/file"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadLines(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		content   string
		wantLines []string
		wantErr   bool
	}{
		{
			name:      "simple lines",
			content:   "line1\nline2\nline3",
			wantLines: []string{"line1", "line2", "line3"},
		},
		{
			name:      "empty file",
			content:   "",
			wantLines: []string{},
		},
		{
			name:      "single line no newline",
			content:   "single line",
			wantLines: []string{"single line"},
		},
		{
			name:      "trailing newline",
			content:   "line1\nline2\n",
			wantLines: []string{"line1", "line2"},
		},
		{
			name:      "empty lines",
			content:   "line1\n\nline3\n",
			wantLines: []string{"line1", "", "line3"},
		},
		{
			name:      "windows line endings",
			content:   "line1\r\nline2\r\nline3",
			wantLines: []string{"line1", "line2", "line3"},
		},
		{
			name:      "unicode content",
			content:   "Hello 世界\n🌍🌎🌏\nTest",
			wantLines: []string{"Hello 世界", "🌍🌎🌏", "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, "test.txt")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Test ReadLines
			lines, err := file.ReadLines(testFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLines, lines)
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		lines, err := file.ReadLines(filepath.Join(tmpDir, "nonexistent.txt"))
		assert.Error(t, err)
		assert.Nil(t, lines)
		assert.Contains(t, err.Error(), "failed to open file")
	})
}

func TestReadContent(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		wantContent string
		wantErr     bool
	}{
		{
			name:        "simple content",
			content:     "Hello, World!",
			wantContent: "Hello, World!",
		},
		{
			name:        "multiline content",
			content:     "Line 1\nLine 2\nLine 3",
			wantContent: "Line 1\nLine 2\nLine 3",
		},
		{
			name:        "empty file",
			content:     "",
			wantContent: "",
		},
		{
			name:        "unicode content",
			content:     "Hello 世界 🌍",
			wantContent: "Hello 世界 🌍",
		},
		{
			name:        "binary content",
			content:     string([]byte{0, 1, 2, 3, 255}),
			wantContent: string([]byte{0, 1, 2, 3, 255}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, "test.txt")
			err := os.WriteFile(testFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			// Test ReadContent
			content, err := file.ReadContent(testFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantContent, content)
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		content, err := file.ReadContent(filepath.Join(tmpDir, "nonexistent.txt"))
		assert.Error(t, err)
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), "failed to read file")
	})
}

func TestReadBytes(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		content   []byte
		wantBytes []byte
		wantErr   bool
	}{
		{
			name:      "text bytes",
			content:   []byte("Hello, World!"),
			wantBytes: []byte("Hello, World!"),
		},
		{
			name:      "binary bytes",
			content:   []byte{0, 1, 2, 3, 127, 128, 255},
			wantBytes: []byte{0, 1, 2, 3, 127, 128, 255},
		},
		{
			name:      "empty bytes",
			content:   []byte{},
			wantBytes: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tmpDir, "test.bin")
			err := os.WriteFile(testFile, tt.content, 0644)
			require.NoError(t, err)

			// Test ReadBytes
			data, err := file.ReadBytes(testFile)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBytes, data)
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		data, err := file.ReadBytes(filepath.Join(tmpDir, "nonexistent.bin"))
		assert.Error(t, err)
		assert.Nil(t, data)
	})
}

func TestStreamLines(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("collect all lines", func(t *testing.T) {
		content := "line1\nline2\nline3\n"
		testFile := filepath.Join(tmpDir, "stream.txt")
		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		var lines []string
		err = file.StreamLines(testFile, func(line string) error {
			lines = append(lines, line)
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"line1", "line2", "line3"}, lines)
	})

	t.Run("stop on error", func(t *testing.T) {
		content := "line1\nline2\nline3\n"
		testFile := filepath.Join(tmpDir, "stream_error.txt")
		err := os.WriteFile(testFile, []byte(content), 0644)
		require.NoError(t, err)

		var lines []string
		stopErr := assert.AnError
		err = file.StreamLines(testFile, func(line string) error {
			lines = append(lines, line)
			if line == "line2" {
				return stopErr
			}
			return nil
		})

		assert.Error(t, err)
		assert.Equal(t, stopErr, err)
		assert.Equal(t, []string{"line1", "line2"}, lines) // Should stop after line2
	})

	t.Run("process large file", func(t *testing.T) {
		// Create a file with many lines
		var content strings.Builder
		lineCount := 10000
		for i := 0; i < lineCount; i++ {
			content.WriteString("line ")
			content.WriteString(string(rune(i)))
			content.WriteByte('\n')
		}

		testFile := filepath.Join(tmpDir, "large.txt")
		err := os.WriteFile(testFile, []byte(content.String()), 0644)
		require.NoError(t, err)

		count := 0
		err = file.StreamLines(testFile, func(line string) error {
			count++
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, lineCount, count)
	})

	t.Run("non-existent file", func(t *testing.T) {
		err := file.StreamLines(filepath.Join(tmpDir, "nonexistent.txt"), func(line string) error {
			t.Fatal("callback should not be called")
			return nil
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open file")
	})
}

func TestReadFrom(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantContent string
		wantErr     bool
	}{
		{
			name:        "simple string",
			input:       "Hello, World!",
			wantContent: "Hello, World!",
		},
		{
			name:        "multiline",
			input:       "Line 1\nLine 2\n",
			wantContent: "Line 1\nLine 2\n",
		},
		{
			name:        "empty",
			input:       "",
			wantContent: "",
		},
		{
			name:        "unicode",
			input:       "Hello 世界 🌍",
			wantContent: "Hello 世界 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			content, err := file.ReadFrom(reader)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantContent, content)
			}
		})
	}

	t.Run("read from bytes.Buffer", func(t *testing.T) {
		buf := bytes.NewBufferString("buffer content")
		content, err := file.ReadFrom(buf)
		assert.NoError(t, err)
		assert.Equal(t, "buffer content", content)
	})
}

func TestFilePermissions(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()

	t.Run("no read permission", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "no_read.txt")
		err := os.WriteFile(testFile, []byte("test"), 0000)
		require.NoError(t, err)

		// ReadLines should fail
		_, err = file.ReadLines(testFile)
		assert.Error(t, err)

		// ReadContent should fail
		_, err = file.ReadContent(testFile)
		assert.Error(t, err)

		// ReadBytes should fail
		_, err = file.ReadBytes(testFile)
		assert.Error(t, err)

		// StreamLines should fail
		err = file.StreamLines(testFile, func(line string) error {
			return nil
		})
		assert.Error(t, err)
	})
}

// Benchmarks
func BenchmarkReadLines(b *testing.B) {
	// Create test file with 1000 lines
	tmpDir := b.TempDir()
	var content strings.Builder
	for i := 0; i < 1000; i++ {
		content.WriteString("This is line number ")
		content.WriteString(string(rune(i)))
		content.WriteByte('\n')
	}

	testFile := filepath.Join(tmpDir, "bench.txt")
	err := os.WriteFile(testFile, []byte(content.String()), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := file.ReadLines(testFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadContent(b *testing.B) {
	// Create test file
	tmpDir := b.TempDir()
	content := strings.Repeat("This is a test line.\n", 1000)
	testFile := filepath.Join(tmpDir, "bench.txt")
	err := os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := file.ReadContent(testFile)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStreamLines(b *testing.B) {
	// Create test file
	tmpDir := b.TempDir()
	var content strings.Builder
	for i := 0; i < 1000; i++ {
		content.WriteString("This is line number ")
		content.WriteString(string(rune(i)))
		content.WriteByte('\n')
	}

	testFile := filepath.Join(tmpDir, "bench.txt")
	err := os.WriteFile(testFile, []byte(content.String()), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count := 0
		err := file.StreamLines(testFile, func(line string) error {
			count++
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
