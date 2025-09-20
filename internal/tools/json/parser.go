package json

import (
	"fmt"
	"strconv"
)

// Parser implements a tokenizer-based JSON parser
type Parser struct {
	tokenizer    *Tokenizer
	currentToken *Token
	peekToken    *Token
}

// NewParser creates a new parser with tokenizer
func NewParser(input string) *Parser {
	p := &Parser{
		tokenizer: NewTokenizer(input),
	}
	// Load first two tokens : after 2nd I'll have current token and next Peek token ready
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances to the next token
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	token, err := p.tokenizer.NextToken()
	if err != nil {
		// Store error token
		p.peekToken = &Token{Type: TOKEN_INVALID, Position: p.tokenizer.position}
	} else {
		p.peekToken = token
	}
}

// Parse validates and parses the JSON input
func (p *Parser) Parse() (interface{}, error) {
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// Ensure we've consumed all input
	if p.currentToken.Type != TOKEN_EOF {
		return nil, p.error("unexpected token after JSON value: %s", TokenTypeName(p.currentToken.Type))
	}

	return value, nil
}

// Validate just checks if JSON is valid (no value building)
func (p *Parser) Validate() error {
	_, err := p.Parse()
	return err
}

// parseValue parses any JSON value
func (p *Parser) parseValue() (interface{}, error) {
	switch p.currentToken.Type {
	case TOKEN_LBRACE:
		return p.parseObject()
	case TOKEN_LBRACKET:
		return p.parseArray()
	case TOKEN_STRING:
		value := p.currentToken.Value
		p.nextToken()
		return value, nil
	case TOKEN_NUMBER:
		value, err := p.parseNumber()
		if err != nil {
			return nil, err
		}
		p.nextToken()
		return value, nil
	case TOKEN_TRUE:
		p.nextToken()
		return true, nil
	case TOKEN_FALSE:
		p.nextToken()
		return false, nil
	case TOKEN_NULL:
		p.nextToken()
		return nil, nil
	case TOKEN_EOF:
		return nil, p.error("unexpected end of JSON input")
	default:
		return nil, p.error("unexpected token: %s", TokenTypeName(p.currentToken.Type))
	}
}

// expectToken checks if current token matches expected type
func (p *Parser) expectToken(tokenType TokenType) error {
	if p.currentToken.Type != tokenType {
		return p.error("expected %s, got %s", TokenTypeName(tokenType), TokenTypeName(p.currentToken.Type))
	}
	return nil
}

// parseObject parses a JSON object
func (p *Parser) parseObject() (map[string]interface{}, error) {
	object := make(map[string]interface{})

	// Consume '{'
	if err := p.expectToken(TOKEN_LBRACE); err != nil {
		return nil, err
	}
	p.nextToken()

	// Handle empty object
	if p.currentToken.Type == TOKEN_RBRACE {
		p.nextToken()
		return object, nil
	}

	for {
		// Parse key
		if err := p.expectToken(TOKEN_STRING); err != nil {
			return nil, p.error("object key must be string")
		}
		key := p.currentToken.Value
		p.nextToken()

		// Check for duplicate key
		if _, exists := object[key]; exists {
			// Note: Standard JSON allows duplicate keys, last one wins
			// Uncomment below to make it strict
			// return nil, p.error("duplicate object key: %s", key)
		}

		// Expect colon
		if err := p.expectToken(TOKEN_COLON); err != nil {
			return nil, err
		}
		p.nextToken()

		// Parse value
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		object[key] = value

		// Check for comma or end
		if p.currentToken.Type == TOKEN_RBRACE {
			p.nextToken()
			return object, nil
		}

		if p.currentToken.Type != TOKEN_COMMA {
			return nil, p.error("expected ',' or '}' after object member")
		}
		p.nextToken()

		// Check for trailing comma
		if p.currentToken.Type == TOKEN_RBRACE {
			return nil, p.error("trailing comma in object")
		}
	}
}

// parseArray parses a JSON array
func (p *Parser) parseArray() ([]interface{}, error) {
	array := make([]interface{}, 0)

	// Consume '['
	if err := p.expectToken(TOKEN_LBRACKET); err != nil {
		return nil, err
	}
	p.nextToken()

	// Handle empty array
	if p.currentToken.Type == TOKEN_RBRACKET {
		p.nextToken()
		return array, nil
	}

	for {
		// Parse value
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		array = append(array, value)

		// Check for comma or end
		if p.currentToken.Type == TOKEN_RBRACKET {
			p.nextToken()
			return array, nil
		}

		if p.currentToken.Type != TOKEN_COMMA {
			return nil, p.error("expected ',' or ']' after array element")
		}
		p.nextToken()

		// Check for trailing comma
		if p.currentToken.Type == TOKEN_RBRACKET {
			return nil, p.error("trailing comma in array")
		}
	}
}

// parseNumber converts number token to appropriate type
func (p *Parser) parseNumber() (interface{}, error) {
	numStr := p.currentToken.Value

	// Try integer first
	if intVal, err := strconv.ParseInt(numStr, 10, 64); err == nil {
		// Check if it fits in int
		if intVal >= -2147483648 && intVal <= 2147483647 {
			return int(intVal), nil
		}
		return intVal, nil
	}

	// Try float
	if floatVal, err := strconv.ParseFloat(numStr, 64); err == nil {
		return floatVal, nil
	}

	return nil, p.error("invalid number: %s", numStr)
}

// error creates an error with token position
func (p *Parser) error(format string, args ...interface{}) error {
	if p.currentToken != nil {
		return fmt.Errorf("JSON parse error at line %d, column %d: %s",
			p.currentToken.Line, p.currentToken.Column, fmt.Sprintf(format, args...))
	}
	return fmt.Errorf("JSON parse error: %s", fmt.Sprintf(format, args...))
}

// ValidateJSON validates JSON using the tokenizer-based parser
func ValidateJSON(input string) error {
	parser := NewParser(input)
	return parser.Validate()
}

// ParseJSON parses JSON and returns the value
func ParseJSON(input string) (interface{}, error) {
	parser := NewParser(input)
	return parser.Parse()
}
