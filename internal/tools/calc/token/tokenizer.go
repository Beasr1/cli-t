package token

import (
	"fmt"
	"strconv"
	"strings"
)

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func extractAndResetNumber(number *strings.Builder) (Token, error) {
	num := strings.Trim(number.String(), " ")
	var token Token
	if len(num) != 0 {
		if !isNumber(num) {
			return Token{}, fmt.Errorf("invalid expression : %s is not a number", num)
		}
		token = Token{
			Value: num,
			Type:  NUMBER,
		}
	} else {
		token = Token{
			Value: "",
			Type:  UNKNOWN,
		}
	}
	number.Reset()
	return token, nil
}

// -3 + 5
func Tokenizer(expression string) ([]Token, error) {
	tokens := make([]Token, 0)

	var number strings.Builder
	for i := 0; i < len(expression); i++ {
		if strings.Contains(OPERATOR_VALUES, string(expression[i])) {
			numberToken, err := extractAndResetNumber(&number)
			if err != nil {
				return []Token{}, err
			}
			if numberToken.Type == NUMBER {
				tokens = append(tokens, numberToken)
			}
			tokens = append(tokens, Token{
				Value: string(expression[i]),
				Type:  OPERATOR,
			})
		} else if strings.Contains(LEFT_PAREN_VALUES, string(expression[i])) {
			numberToken, err := extractAndResetNumber(&number)
			if err != nil {
				return []Token{}, err
			}
			if numberToken.Type == NUMBER {
				tokens = append(tokens, numberToken)
			}
			tokens = append(tokens, Token{
				Value: string(expression[i]),
				Type:  LEFT_PAREN,
			})
		} else if strings.Contains(RIGHT_PAREN_VALUES, string(expression[i])) {
			numberToken, err := extractAndResetNumber(&number)
			if err != nil {
				return []Token{}, err
			}
			if numberToken.Type == NUMBER {
				tokens = append(tokens, numberToken)
			}
			tokens = append(tokens, Token{
				Value: string(expression[i]),
				Type:  RIGHT_PAREN,
			})
		} else {
			number.WriteByte(expression[i])
		}
	}

	numberToken, err := extractAndResetNumber(&number)
	if err != nil {
		return []Token{}, err
	}
	if numberToken.Type == NUMBER {
		tokens = append(tokens, numberToken)
	}

	return tokens, nil
}
