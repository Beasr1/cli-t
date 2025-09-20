package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test internal helper functions
func TestTokenizer_HelperFunctions(t *testing.T) {
	// Test isDigit
	assert.True(t, isDigit('0'))
	assert.True(t, isDigit('5'))
	assert.True(t, isDigit('9'))
	assert.False(t, isDigit('a'))
	assert.False(t, isDigit(' '))
	assert.False(t, isDigit('.'))

	// Test isLetter
	assert.True(t, isLetter('a'))
	assert.True(t, isLetter('z'))
	assert.True(t, isLetter('A'))
	assert.True(t, isLetter('Z'))
	assert.False(t, isLetter('0'))
	assert.False(t, isLetter(' '))
	assert.False(t, isLetter('_'))
}

// Test TokenTypeName function
func TestTokenTypeName(t *testing.T) {
	tests := []struct {
		token    TokenType
		expected string
	}{
		{TOKEN_EOF, "EOF"},
		{TOKEN_LBRACE, "LBRACE"},
		{TOKEN_RBRACE, "RBRACE"},
		{TOKEN_LBRACKET, "LBRACKET"},
		{TOKEN_RBRACKET, "RBRACKET"},
		{TOKEN_COLON, "COLON"},
		{TOKEN_COMMA, "COMMA"},
		{TOKEN_STRING, "STRING"},
		{TOKEN_NUMBER, "NUMBER"},
		{TOKEN_TRUE, "TRUE"},
		{TOKEN_FALSE, "FALSE"},
		{TOKEN_NULL, "NULL"},
		{TOKEN_INVALID, "INVALID"},
		{TokenType(999), "UNKNOWN"}, // Out of range
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, TokenTypeName(tt.token))
		})
	}
}

// Test internal methods
func TestTokenizer_current(t *testing.T) {
	tok := &Tokenizer{
		input:    "test",
		position: 0,
	}

	// Test current at various positions
	assert.Equal(t, byte('t'), tok.current())

	tok.position = 3
	assert.Equal(t, byte('t'), tok.current())

	tok.position = 4 // Past end
	assert.Equal(t, byte(0), tok.current())
}

func TestTokenizer_peek(t *testing.T) {
	tok := &Tokenizer{
		input:    "test",
		position: 0,
	}

	// Test peek at various positions
	assert.Equal(t, byte('e'), tok.peek())

	tok.position = 3                     // Last character
	assert.Equal(t, byte(0), tok.peek()) // Nothing to peek
}

func TestTokenizer_advance(t *testing.T) {
	tok := &Tokenizer{
		input:    "a\nb",
		position: 0,
		line:     1,
		column:   1,
	}

	// Advance through normal character
	assert.Equal(t, 0, tok.position)
	assert.Equal(t, 1, tok.line)
	assert.Equal(t, 1, tok.column)

	tok.advance()
	assert.Equal(t, 1, tok.position)
	assert.Equal(t, 1, tok.line)
	assert.Equal(t, 2, tok.column)

	// Advance through newline
	tok.advance()
	assert.Equal(t, 2, tok.position)
	assert.Equal(t, 2, tok.line)
	assert.Equal(t, 1, tok.column)

	// Advance at end
	tok.position = len(tok.input)
	tok.advance() // Should not panic
	assert.Equal(t, len(tok.input), tok.position)
}

// Test error method
func TestTokenizer_error(t *testing.T) {
	tok := &Tokenizer{
		input:    "test",
		position: 5,
		line:     2,
		column:   3,
	}

	err := tok.error("test error: %s", "details")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line 2, column 3")
	assert.Contains(t, err.Error(), "test error: details")
}

// Test skipWhitespace with various inputs
func TestTokenizer_skipWhitespace(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		startPos     int
		expectedPos  int
		expectedLine int
		expectedCol  int
	}{
		{
			name:         "no whitespace",
			input:        "abc",
			startPos:     0,
			expectedPos:  0,
			expectedLine: 1,
			expectedCol:  1,
		},
		{
			name:         "spaces",
			input:        "   abc",
			startPos:     0,
			expectedPos:  3,
			expectedLine: 1,
			expectedCol:  4,
		},
		{
			name:         "tabs",
			input:        "\t\tabc",
			startPos:     0,
			expectedPos:  2,
			expectedLine: 1,
			expectedCol:  3,
		},
		{
			name:         "newlines",
			input:        "\n\nabc",
			startPos:     0,
			expectedPos:  2,
			expectedLine: 3,
			expectedCol:  1,
		},
		{
			name:         "mixed whitespace",
			input:        " \t\n \r abc",
			startPos:     0,
			expectedPos:  6,
			expectedLine: 2,
			expectedCol:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := &Tokenizer{
				input:    tt.input,
				position: tt.startPos,
				line:     1,
				column:   1,
			}

			tok.skipWhitespace()

			assert.Equal(t, tt.expectedPos, tok.position)
			assert.Equal(t, tt.expectedLine, tok.line)
			assert.Equal(t, tt.expectedCol, tok.column)
		})
	}
}

// Test readKeyword edge cases
func TestTokenizer_readKeyword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "true keyword",
			input:    "true",
			expected: "true",
		},
		{
			name:     "false keyword",
			input:    "false",
			expected: "false",
		},
		{
			name:     "null keyword",
			input:    "null",
			expected: "null",
		},
		{
			name:     "partial keyword",
			input:    "tru",
			expected: "tru",
		},
		{
			name:     "keyword with trailing",
			input:    "true123",
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := &Tokenizer{
				input:    tt.input,
				position: 0,
			}

			result, err := tok.readKeyword()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test concurrent tokenizer usage
func TestTokenizer_Concurrent(t *testing.T) {
	// Each goroutine gets its own tokenizer
	input := `{"test": true}`

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			tok := NewTokenizer(input)
			token, err := tok.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, TOKEN_LBRACE, token.Type)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Benchmark tokenizer initialization
func BenchmarkNewTokenizer(b *testing.B) {
	input := `{"test": true}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewTokenizer(input)
	}
}

// Benchmark skipWhitespace
func BenchmarkTokenizer_skipWhitespace(b *testing.B) {
	input := "    \t\n    \r\n    {}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tok := &Tokenizer{
			input:    input,
			position: 0,
			line:     1,
			column:   1,
		}
		tok.skipWhitespace()
	}
}
