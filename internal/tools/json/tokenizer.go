package json

import (
	"fmt"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// TokenType represents the type of JSON token
type TokenType int

const (
	TOKEN_EOF      TokenType = iota
	TOKEN_LBRACE             // {
	TOKEN_RBRACE             // }
	TOKEN_LBRACKET           // [
	TOKEN_RBRACKET           // ]
	TOKEN_COLON              // :
	TOKEN_COMMA              // ,
	TOKEN_STRING             // "xxx"
	TOKEN_NUMBER             // 1
	TOKEN_TRUE               // true
	TOKEN_FALSE              // false
	TOKEN_NULL               // null
	TOKEN_INVALID            // atmkbfg
)

// Token represents a JSON token with its type, value, and position
type Token struct {
	Type     TokenType
	Value    string
	Position int
	Line     int
	Column   int
}

// Tokenizer performs lexical analysis of JSON input
type Tokenizer struct {
	input    string
	position int
	line     int
	column   int
}

// NewTokenizer creates a new tokenizer
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
	}
}

// NextToken returns the next token from input
func (t *Tokenizer) NextToken() (*Token, error) {
	// Skip whitespace
	t.skipWhitespace()

	// Check EOF or check if t.current() is 0
	if t.position >= len(t.input) {
		return &Token{Type: TOKEN_EOF, Position: t.position, Line: t.line, Column: t.column}, nil
	}

	// Save position for token
	tokenPos := t.position
	tokenLine := t.line
	tokenCol := t.column

	ch := t.current()

	switch ch {
	case '{':
		t.advance()
		return &Token{Type: TOKEN_LBRACE, Value: "{", Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case '}':
		t.advance()
		return &Token{Type: TOKEN_RBRACE, Value: "}", Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case '[':
		t.advance()
		return &Token{Type: TOKEN_LBRACKET, Value: "[", Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case ']':
		t.advance()
		return &Token{Type: TOKEN_RBRACKET, Value: "]", Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case ':':
		t.advance()
		return &Token{Type: TOKEN_COLON, Value: ":", Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case ',':
		t.advance()
		return &Token{Type: TOKEN_COMMA, Value: ",", Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case '"':
		value, err := t.readString()
		if err != nil {
			return nil, err
		}
		return &Token{Type: TOKEN_STRING, Value: value, Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case 't', 'f':
		keyword, err := t.readKeyword()
		if err != nil {
			return nil, err
		}
		switch keyword {
		case "true":
			return &Token{Type: TOKEN_TRUE, Value: keyword, Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
		case "false":
			return &Token{Type: TOKEN_FALSE, Value: keyword, Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
		default:
			return nil, t.error("invalid keyword: %s", keyword)
		}
	case 'n':
		keyword, err := t.readKeyword()
		if err != nil {
			return nil, err
		}
		if keyword != "null" {
			return nil, t.error("invalid keyword: %s", keyword)
		}
		return &Token{Type: TOKEN_NULL, Value: keyword, Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		value, err := t.readNumber()
		if err != nil {
			return nil, err
		}
		return &Token{Type: TOKEN_NUMBER, Value: value, Position: tokenPos, Line: tokenLine, Column: tokenCol}, nil
	default:
		return nil, t.error("unexpected character '%c'", ch)
	}
}

// Helper methods

func (t *Tokenizer) current() byte {
	if t.position >= len(t.input) {
		return 0 // 0 value if iota for EOF
	}
	return t.input[t.position]
}

func (t *Tokenizer) peek() byte {
	if t.position+1 >= len(t.input) {
		return 0
	}
	return t.input[t.position+1]
}

func (t *Tokenizer) advance() {
	if t.position < len(t.input) {
		if t.input[t.position] == '\n' {
			t.line++
			t.column = 1
		} else {
			t.column++
		}
		t.position++
	}
}

// Keeps on skipping whitespaces
func (t *Tokenizer) skipWhitespace() {
	for t.position < len(t.input) {
		ch := t.current()
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			t.advance()
		} else {
			break
		}
	}
}

func (t *Tokenizer) readString() (string, error) {
	t.advance() // Skip opening quote

	// += Allocates new memory for the combined string : Copies all previous characters O(N2)
	var result strings.Builder // better efficiency as it actually modify the same string and not create a new string

	for t.position < len(t.input) {
		ch := t.current()

		if ch == '"' {
			t.advance()                 // Skip closing quote
			return result.String(), nil // will only create string from the buffer once (better efficiency)
		}

		// checks whether character is \
		if ch == '\\' {
			t.advance()
			if t.position >= len(t.input) {
				return "", t.error("unterminated string")
			}

			switch t.current() {
			case '"':
				result.WriteByte('"')
			case '\\':
				result.WriteByte('\\')
			case '/':
				result.WriteByte('/')
			case 'b':
				result.WriteByte('\b')
			case 'f':
				result.WriteByte('\f')
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			case 'u':
				t.advance()
				// \u0041 ... etc
				r, err := t.readUnicodeEscape()
				if err != nil {
					return "", err
				}
				result.WriteRune(r)
				continue // Don't advance again.
			default:
				return "", t.error("invalid escape sequence '\\%c'", t.current())
			}
		} else if ch < 0x20 {
			return "", t.error("unescaped control character in string")
		} else {
			result.WriteByte(ch)
		}

		t.advance()
	}

	return "", t.error("unterminated string")
}

// Basic Multilingual Plane (BMP)
func (t *Tokenizer) readUnicodeEscape() (rune, error) {
	if t.position+3 >= len(t.input) {
		return 0, t.error("incomplete unicode escape sequence")
	}

	hex := t.input[t.position : t.position+4] // UTF-16 = 4 bytes. json uses
	var code uint16

	for i := 0; i < 4; i++ {
		digit := hex[i]
		code <<= 4

		switch {
		case digit >= '0' && digit <= '9':
			code |= uint16(digit - '0')
		case digit >= 'a' && digit <= 'f':
			code |= uint16(digit - 'a' + 10)
		case digit >= 'A' && digit <= 'F': // case insensitive.
			code |= uint16(digit - 'A' + 10)
		default:
			return 0, t.error("invalid unicode escape sequence")
		}
	}

	t.position += 4
	t.column += 4

	// Handle surrogate pairs if needed
	// Unicode reserves special ranges:
	// High Surrogates: 0xD800 to 0xDBFF (1,024 values)
	// Low Surrogates:  0xDC00 to 0xDFFF (1,024 values)
	if utf16.IsSurrogate(rune(code)) {
		// For simplicity, we'll just return the replacement character
		// A full implementation would handle surrogate pairs
		return utf8.RuneError, nil
	}

	return rune(code), nil
}

func (t *Tokenizer) readNumber() (string, error) {
	start := t.position

	// Optional minus
	if t.current() == '-' {
		t.advance()
		if t.position >= len(t.input) || !isDigit(t.current()) {
			return "", t.error("invalid number")
		}
	}

	// Integer part
	if t.current() == '0' {
		t.advance()
		// After 0, only . or e/E or end is valid
		if t.position < len(t.input) && isDigit(t.current()) {
			return "", t.error("invalid number: leading zeros not allowed")
		}
	} else {
		// Must start with 1-9
		if !isDigit(t.current()) {
			return "", t.error("invalid number")
		}
		for t.position < len(t.input) && isDigit(t.current()) {
			t.advance()
		}
	}

	// Fractional part
	if t.position < len(t.input) && t.current() == '.' {
		t.advance()
		if t.position >= len(t.input) || !isDigit(t.current()) {
			return "", t.error("invalid number: expected digit after decimal point")
		}
		for t.position < len(t.input) && isDigit(t.current()) {
			t.advance()
		}
	}

	// Exponent part
	if t.position < len(t.input) && (t.current() == 'e' || t.current() == 'E') {
		t.advance()
		if t.position < len(t.input) && (t.current() == '+' || t.current() == '-') {
			t.advance()
		}
		if t.position >= len(t.input) || !isDigit(t.current()) {
			return "", t.error("invalid number: expected digit in exponent")
		}
		for t.position < len(t.input) && isDigit(t.current()) {
			t.advance()
		}
	}

	return t.input[start:t.position], nil
}

func (t *Tokenizer) readKeyword() (string, error) {
	start := t.position

	for t.position < len(t.input) && isLetter(t.current()) {
		t.advance()
	}

	return t.input[start:t.position], nil
}

func (t *Tokenizer) error(format string, args ...interface{}) error {
	return fmt.Errorf("JSON parse error at line %d, column %d: %s",
		t.line, t.column, fmt.Sprintf(format, args...))
}

// Helper functions
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// TokenTypeName returns the string name of a token type
func TokenTypeName(t TokenType) string {
	names := []string{
		"EOF",
		"LBRACE",
		"RBRACE",
		"LBRACKET",
		"RBRACKET",
		"COLON",
		"COMMA",
		"STRING",
		"NUMBER",
		"TRUE",
		"FALSE",
		"NULL",
		"INVALID",
	}
	if int(t) < len(names) {
		return names[t]
	}
	return "UNKNOWN"
}
