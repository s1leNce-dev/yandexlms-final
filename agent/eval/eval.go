package eval

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Eval(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	var numStack []float64
	var opStack []rune

	priority := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	applyOperator := func() error {
		if len(numStack) < 2 || len(opStack) == 0 {
			return fmt.Errorf("ошибка в выражении")
		}

		b, a := numStack[len(numStack)-1], numStack[len(numStack)-2]
		numStack = numStack[:len(numStack)-2]

		op := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]

		var result float64
		switch op {
		case '+':
			result = a + b
		case '-':
			result = a - b
		case '*':
			result = a * b
		case '/':
			if b == 0 {
				return fmt.Errorf("деление на ноль")
			}
			result = a / b
		}

		numStack = append(numStack, result)
		return nil
	}

	for i := 0; i < len(expression); i++ {
		ch := rune(expression[i])

		if unicode.IsDigit(ch) || ch == '.' {
			start := i
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				i++
			}
			value, err := strconv.ParseFloat(expression[start:i], 64)
			if err != nil {
				return 0, fmt.Errorf("неверный формат числа")
			}
			numStack = append(numStack, value)
			i--
		} else if ch == '(' {
			opStack = append(opStack, ch)
		} else if ch == ')' {
			for len(opStack) > 0 && opStack[len(opStack)-1] != '(' {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			if len(opStack) == 0 || opStack[len(opStack)-1] != '(' {
				return 0, fmt.Errorf("ошибка в расстановке скобок")
			}
			opStack = opStack[:len(opStack)-1]
		} else if strings.ContainsRune("+-*/", ch) {
			for len(opStack) > 0 && opStack[len(opStack)-1] != '(' &&
				priority[opStack[len(opStack)-1]] >= priority[ch] {
				if err := applyOperator(); err != nil {
					return 0, err
				}
			}
			opStack = append(opStack, ch)
		} else {
			return 0, fmt.Errorf("недопустимый символ: %c", ch)
		}
	}

	for len(opStack) > 0 {
		if opStack[len(opStack)-1] == '(' {
			return 0, fmt.Errorf("ошибка в расстановке скобок")
		}
		if err := applyOperator(); err != nil {
			return 0, err
		}
	}

	if len(numStack) != 1 {
		return 0, fmt.Errorf("ошибка в выражении")
	}

	return numStack[0], nil
}
