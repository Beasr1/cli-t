package stack

import (
	"cli-t/internal/tools/calc/token"
	"fmt"
	"math"
)

// Precedence levels
const (
	LOWEST   = iota
	SUM      // + -
	PRODUCT  // * / %
	EXPONENT // ^
)

func GetPrecedence(op string) int {
	switch op {
	case "+", "-":
		return SUM
	case "*", "/", "%":
		return PRODUCT
	case "^":
		return EXPONENT
	default:
		return LOWEST
	}
}

func IsOperator(t token.Token) bool {
	return t.Type == token.OPERATOR
}

// Right-associative: ^ (power)
// Left-associative: everything else
func IsRightAssociative(op string) bool {
	return op == "^"
}

func ApplyOperator(op string, left, right float64) (float64, error) {
	switch op {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	case "%":
		if right == 0 {
			return 0, fmt.Errorf("modulo by zero")
		}
		return math.Mod(left, right), nil
	case "^":
		return math.Pow(left, right), nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}
