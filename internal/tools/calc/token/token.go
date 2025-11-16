package token

type TokenType int

const (
	UNKNOWN TokenType = iota
	NUMBER
	OPERATOR
	LEFT_PAREN
	RIGHT_PAREN
)

type Token struct {
	Type  TokenType
	Value string // "3.14", "+", "(", etc.
}

// TODO: should we allow other brackets
const OPERATOR_VALUES = "+-/*%"
const LEFT_PAREN_VALUES = "("
const RIGHT_PAREN_VALUES = ")"
