package calculator_test

import (
	"testing"

	"github.com/wifi538/CalculatorOnline/pkg/calculator"
)

// основная функция тестирования
func TestCalc(t *testing.T) {
	//тесты с успешным выполнением
	testCasesSuccess := []struct {
		name           string  //название теста
		expression     string  //выражение
		expectedResult float64 //ожидаемый результат
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(1+2)*2-3*(3)",
			expectedResult: -3,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	//проверка каждого случая из testCasesSuccess
	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			//вызываем Calc с тестовым значением
			val, err := calculator.Calc(testCase.expression)
			//при ошибке тест завершается с ошибкой
			if err != nil {
				t.Fatalf("successful case %s returns error", testCase.expression)
			}
			//проверка результата с ожиданием
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	//тесты с ошибками
	testCasesFail := []struct {
		name        string //название теста
		expression  string //выражение
		expectedErr error  //ожидаемая ошибка
	}{
		{
			name:       "last is operator",
			expression: "1+1*",
		},
		{
			name:       "two operators together",
			expression: "1+2**2",
		},
		{
			name:       "opened and not closed bracket",
			expression: "(2+2",
		},
		{
			name:       "division by zero",
			expression: "2 / (1 - 1)",
		},
		{
			name:       "wrong character",
			expression: "2 + 1a",
		},
		{
			name:       "no symbol between brackets",
			expression: "(2-1)(1+2)",
		},
		{
			name:       "empty",
			expression: "",
		},
	}

	//проверка каждого случая из testCasesFail
	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			//вызываем Calc с тестовым значением
			val, err := calculator.Calc(testCase.expression)
			//если ошибки не было, тест завершается с ошибкой
			if err == nil {
				t.Fatalf("expression %s is invalid but result %f was obtained", testCase.expression, val)
			}
		})
	}
}
