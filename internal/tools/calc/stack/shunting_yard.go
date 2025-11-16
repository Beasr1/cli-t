package stack

import (
	"cli-t/internal/tools/calc/token"
	"fmt"
)

/*
3-2^2 =  322^-
3-2+2 =  32-2+ left ass
3*2/3 =  32*3/ left ass
3^2^2 =  322^^ right ass

3-(3+2) = 332+-


Shunting Yard

Converting infix notation to postfix (RPN)
Understanding Reverse Polish Notation
Why postfix is easier to evaluate
*/

func ToPostfix(tokens []token.Token) ([]token.Token, error) {
	output := make([]token.Token, 0)
	opStack := make([]token.Token, 0)

	for _, t := range tokens {
		if t.Type == token.NUMBER {
			output = append(output, t)
		}

		if t.Type == token.OPERATOR {
			for {

				// break if empty or not operator (parenthesis reached)
				opLen := len(opStack)
				if opLen == 0 || opStack[opLen-1].Type != token.OPERATOR {
					break
				}

				leftAss := !IsRightAssociative(opStack[opLen-1].Value)
				if leftAss {
					if GetPrecedence(opStack[opLen-1].Value) >= GetPrecedence(t.Value) {
						// if equal than also pop stack top and put that in output : since left ass
						output = append(output, opStack[opLen-1])
						opStack = opStack[:opLen-1] // start , size
					} else {
						break
					}
				} else {
					if GetPrecedence(opStack[opLen-1].Value) > GetPrecedence(t.Value) {
						// if less than pop stack top
						output = append(output, opStack[opLen-1])
						opStack = opStack[:opLen-1] // start , size
					} else {
						break
					}
				}
			}
			opStack = append(opStack, t)
		}

		if t.Type == token.LEFT_PAREN {
			opStack = append(opStack, t)
		}

		if t.Type == token.RIGHT_PAREN {
			for {
				// break if empty or not operator (parenthesis reached)
				opLen := len(opStack)
				if opLen == 0 {
					return []token.Token{}, fmt.Errorf("invalid expression")
				}

				if opStack[opLen-1].Type != token.LEFT_PAREN {
					// pop operator
					output = append(output, opStack[opLen-1])
					opStack = opStack[:opLen-1] // start , size
				} else {
					opStack = opStack[:opLen-1] // start , size
					break
				}
			}
		}
	}

	for {
		opLen := len(opStack)
		if opLen == 0 {
			break
		}

		if opStack[opLen-1].Type == token.LEFT_PAREN {
			return []token.Token{}, fmt.Errorf("invalid expression")
		}

		// right paren here is not possible
		output = append(output, opStack[opLen-1])
		opStack = opStack[:opLen-1] // start , size
	}

	return output, nil
}
