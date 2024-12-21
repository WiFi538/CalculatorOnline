package calculator

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func Calc(expression string) (float64, error) {
	numStack := make([]float64, 0)
	opStack := make([]rune, 0)

	doOp := func(op rune) error {
		if len(numStack) < 2 {
			return errors.New("недостаточно операндов для операции")
		}
		b := numStack[len(numStack)-1]
		a := numStack[len(numStack)-2]
		numStack = numStack[:len(numStack)-2]
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
				return errors.New("деление на ноль")
			}
			result = a / b
		default:
			return errors.New("неизвестная операция")
		}
		numStack = append(numStack, result)
		return nil
	}

	//определение приоритета операции
	priority := func(op rune) int {
		switch op {
		case '+', '-':
			return 1
		case '*', '/':
			return 2
		}
		return 0
	}

	//обработка чисел
	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])
		if unicode.IsDigit(char) {
			num, err := strconv.ParseFloat(string(char), 64)
			if err != nil {
				return 0, err
			}
			numStack = append(numStack, num)

			//обработка открывающей и закрывающей скобок
		} else if char == '(' {
			opStack = append(opStack, char)
		} else if char == ')' {
			for len(opStack) > 0 && opStack[len(opStack)-1] != '(' {
				if err := doOp(opStack[len(opStack)-1]); err != nil {
					return 0, err
				}
				opStack = opStack[:len(opStack)-1]
			}
			if len(opStack) == 0 {
				return 0, errors.New("несогласованные скобки")
			}
			opStack = opStack[:len(opStack)-1]

			//обработка операций
		} else if char == '+' || char == '-' || char == '*' || char == '/' {
			for len(opStack) > 0 && priority(opStack[len(opStack)-1]) >= priority(char) {
				if err := doOp(opStack[len(opStack)-1]); err != nil {
					return 0, err
				}
				opStack = opStack[:len(opStack)-1]
			}
			opStack = append(opStack, char)

			//обработка пробелов
		} else if char != ' ' {
			return 0, errors.New("недопустимый символ в выражении")
		}
	}

	for len(opStack) > 0 {
		if opStack[len(opStack)-1] == '(' {
			return 0, errors.New("непарные скобки")
		}
		if err := doOp(opStack[len(opStack)-1]); err != nil {
			return 0, err
		}
		opStack = opStack[:len(opStack)-1]
	}

	//результат один элемент стека
	if len(numStack) != 1 {
		return 0, errors.New("некорректное выражение")
	}

	return numStack[0], nil
}

func main() {
	expression := "3 + 5 * (2 - 8)"
	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Результат:", result)
	}
}
