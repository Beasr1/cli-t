package stack

import (
	"cli-t/internal/tools/calc/token"
	"fmt"
	"strconv"
)

func EvaluatePostfix(tokens []token.Token) (float64, error) {
	var resultToken []float64

	// tokens only have number and operator
	for _, t := range tokens {
		switch t.Type {
		case token.OPERATOR:

			first, second := len(resultToken)-2, len(resultToken)-1
			if first < 0 {
				return 0, fmt.Errorf("invalid expression")
			}

			res, err := ApplyOperator(t.Value, resultToken[first], resultToken[second])
			if err != nil {
				return 0, err
			}

			resultToken = append(resultToken[:first], res)

		case token.NUMBER:

			num, err := strconv.ParseFloat(t.Value, 64)
			if err != nil {
				// already validated number when tokenizing
				return 0, err
			}
			resultToken = append(resultToken, num)

		default:
			return 0, fmt.Errorf("invalid expression, type")
		}
	}

	if len(resultToken) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	return resultToken[0], nil
}
