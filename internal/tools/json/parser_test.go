package json

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test parser creation and initialization
func TestNewParser(t *testing.T) {
	input := `{"test": true}`
	p := NewParser(input)

	assert.NotNil(t, p)
	assert.NotNil(t, p.tokenizer)
	assert.NotNil(t, p.currentToken)
	assert.NotNil(t, p.peekToken)
}

// Test expectToken method
func TestParser_expectToken(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expect  TokenType
		wantErr bool
	}{
		{
			name:    "expect LBRACE success",
			input:   `{"test": true}`,
			expect:  TOKEN_LBRACE,
			wantErr: false,
		},
		{
			name:    "expect LBRACE failure",
			input:   `["test"]`,
			expect:  TOKEN_LBRACE,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			err := p.expectToken(tt.expect)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test parseNumber method
func TestParser_parseNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantType string
	}{
		{
			name:     "small integer",
			input:    "42",
			expected: 42,
			wantType: "int",
		},
		{
			name:     "negative integer",
			input:    "-123",
			expected: -123,
			wantType: "int",
		},
		{
			name:     "float",
			input:    "3.14159",
			expected: 3.14159,
			wantType: "float64",
		},
		{
			name:     "scientific notation",
			input:    "1.23e10",
			expected: 1.23e10,
			wantType: "float64",
		},
		{
			name:     "large int64",
			input:    "9223372036854775807",
			expected: int64(9223372036854775807),
			wantType: "int64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create parser with just the number
			p := NewParser(tt.input)

			// parseNumber expects currentToken to be a number
			assert.Equal(t, TOKEN_NUMBER, p.currentToken.Type)

			result, err := p.parseNumber()
			require.NoError(t, err)

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.wantType, fmt.Sprintf("%T", result))
		})
	}
}

// Test parseObject method with various cases
func TestParser_parseObject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "empty object",
			input:    `{}`,
			expected: map[string]interface{}{},
		},
		{
			name:  "simple object",
			input: `{"key": "value"}`,
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:  "multiple keys",
			input: `{"a": 1, "b": 2, "c": 3}`,
			expected: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		{
			name:  "nested object",
			input: `{"outer": {"inner": "value"}}`,
			expected: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
		},
		{
			name:    "trailing comma",
			input:   `{"a": 1,}`,
			wantErr: true,
			errMsg:  "trailing comma",
		},
		{
			name:    "missing colon",
			input:   `{"key" "value"}`,
			wantErr: true,
			errMsg:  "expected COLON",
		},
		{
			name:    "unquoted key",
			input:   `{key: "value"}`,
			wantErr: true,
			errMsg:  "object key must be string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			result, err := p.parseObject()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Test parseArray method
func TestParser_parseArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []interface{}
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "empty array",
			input:    `[]`,
			expected: []interface{}{},
		},
		{
			name:     "simple array",
			input:    `[1, 2, 3]`,
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "mixed types",
			input:    `[1, "two", true, null]`,
			expected: []interface{}{1, "two", true, nil},
		},
		{
			name:  "nested arrays",
			input: `[[1, 2], [3, 4]]`,
			expected: []interface{}{
				[]interface{}{1, 2},
				[]interface{}{3, 4},
			},
		},
		{
			name:    "trailing comma",
			input:   `[1, 2,]`,
			wantErr: true,
			errMsg:  "trailing comma",
		},
		{
			name:    "missing comma",
			input:   `[1 2 3]`,
			wantErr: true,
			errMsg:  "expected ',' or ']'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			result, err := p.parseArray()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Test parseValue with different token types
func TestParser_parseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantType string
	}{
		{
			name:     "object",
			input:    `{"test": true}`,
			expected: map[string]interface{}{"test": true},
			wantType: "map[string]interface {}",
		},
		{
			name:     "array",
			input:    `[1, 2, 3]`,
			expected: []interface{}{1, 2, 3},
			wantType: "[]interface {}",
		},
		{
			name:     "string",
			input:    `"hello"`,
			expected: "hello",
			wantType: "string",
		},
		{
			name:     "number",
			input:    `42`,
			expected: 42,
			wantType: "int",
		},
		{
			name:     "true",
			input:    `true`,
			expected: true,
			wantType: "bool",
		},
		{
			name:     "false",
			input:    `false`,
			expected: false,
			wantType: "bool",
		},
		{
			name:     "null",
			input:    `null`,
			expected: nil,
			wantType: "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			result, err := p.parseValue()

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)

			actualType := fmt.Sprintf("%T", result)
			assert.Equal(t, tt.wantType, actualType)
		})
	}
}

// Test error handling
func TestParser_error(t *testing.T) {
	p := NewParser(`{"test": true}`)

	// Test with current token
	err := p.error("test error %s", "details")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "line 1, column 1")
	assert.Contains(t, err.Error(), "test error details")
}

// Test token advancement
func TestParser_nextToken(t *testing.T) {
	input := `{"key": "value"}`
	p := NewParser(input)

	// Initial state: currentToken should be first real token
	assert.Equal(t, TOKEN_LBRACE, p.currentToken.Type)

	// Advance
	p.nextToken()
	assert.Equal(t, TOKEN_STRING, p.currentToken.Type)
	assert.Equal(t, "key", p.currentToken.Value)

	// Advance again
	p.nextToken()
	assert.Equal(t, TOKEN_COLON, p.currentToken.Type)

	// Continue advancing
	p.nextToken()
	assert.Equal(t, TOKEN_STRING, p.currentToken.Type)
	assert.Equal(t, "value", p.currentToken.Value)
}

// Test complex nested structures
func TestParser_ComplexStructures(t *testing.T) {
	input := `{
		"name": "test",
		"nested": {
			"array": [1, 2, {"deep": true}],
			"value": null
		},
		"list": [
			{"id": 1},
			{"id": 2}
		]
	}`

	p := NewParser(input)
	result, err := p.Parse()

	require.NoError(t, err)

	obj := result.(map[string]interface{})
	assert.Equal(t, "test", obj["name"])

	nested := obj["nested"].(map[string]interface{})
	array := nested["array"].([]interface{})
	assert.Len(t, array, 3)
	assert.Equal(t, 1, array[0])
	assert.Equal(t, 2, array[1])

	deepObj := array[2].(map[string]interface{})
	assert.Equal(t, true, deepObj["deep"])
}

// Benchmark internal methods
func BenchmarkParser_parseObject(b *testing.B) {
	input := `{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := NewParser(input)
		_, _ = p.parseObject()
	}
}

func BenchmarkParser_parseArray(b *testing.B) {
	input := `[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := NewParser(input)
		_, _ = p.parseArray()
	}
}
