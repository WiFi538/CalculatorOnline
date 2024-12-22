package application_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wifi538/CalculatorOnline/internal/application"
	"github.com/wifi538/CalculatorOnline/pkg/calculator"
)

// структура для десериализации успешного ответа
type ResultResponse struct {
	Result float64 `json:"result"` //результат вычисления
}

// структура для десериализации ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"` //текст ошибки
}

// проверка на корректность обработки запросов
func TestApplication(t *testing.T) {
	//набор тестов с успешной обработкой
	testCasesSuccess := []struct {
		name        string         //название теста
		expression  []byte         //выражение
		expectedRes ResultResponse //ожидаемый ответ
		status      int            //ожидаемый статус HTTP
	}{
		{
			name:        "simple",
			expression:  []byte(`{"expression":"1 + 1"}`),
			expectedRes: ResultResponse{Result: 2},
			status:      http.StatusOK,
		},
		{
			name:        "priority",
			expression:  []byte(`{"expression":"( 2 + 2 ) * 2"}`),
			expectedRes: ResultResponse{Result: 8},
			status:      http.StatusOK,
		},
	}

	//проверка каждого теста из testCasesSuccess
	for _, TestCase := range testCasesSuccess {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(TestCase.expression))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		application.СalcHandler(w, req)
		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		//десериализируем ответ в ResultResponse
		var actualResults ResultResponse
		err = json.Unmarshal(data, &actualResults)
		if err != nil {
			t.Fatal(err)
		}

		//проверка совпадения статуса ответа с ожидаемым ответом
		if TestCase.expectedRes != actualResults {
			t.Fatalf("Test: %s, Expected result: %v, but got: %v", TestCase.name, data, TestCase.expectedRes)
		}
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Test: %s, Expected status: %d, but got: %d", TestCase.name, http.StatusOK, res.StatusCode)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)

	w := httptest.NewRecorder()
	application.СalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	//десериализируем ответ с ошибкой
	var actualErr ErrorResponse
	err = json.Unmarshal(data, &actualErr)
	if err != nil {
		t.Fatal(err)
	}
	//ожидаемый ответ с ошибкой
	expectedErr := ErrorResponse{Error: "invalid request method"}
	if expectedErr != actualErr {
		t.Fatalf("Expected error: %s, but got: %s", expectedErr, actualErr)
	}

	//проверяем совпадение статуса ответа с ожидаемым ответом
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Expected status: %d, but got: %d", http.StatusMethodNotAllowed, res.StatusCode)
	}

	//набор тестов с ошибками
	testCasesFail := []struct {
		name        string        //название теста
		expression  []byte        //выражение
		expectedErr ErrorResponse //ожидаемый ответ с ошибкой
		status      int           //ожидаемый статус HTTP
	}{
		{
			name:        "invalid body",
			expression:  []byte(`aaa`),
			expectedErr: ErrorResponse{Error: "invalid request body"},
			status:      http.StatusMethodNotAllowed,
		},
		{
			name:        "wrong character",
			expression:  []byte(`{"expression":"2 + a"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrWrongCharacter.Error()},
			status:      http.StatusUnprocessableEntity,
		},
		{
			name:        "empty brackets",
			expression:  []byte(`{"expression":"()"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrEmptyBrackets.Error()},
			status:      http.StatusUnprocessableEntity,
		},
		{
			name:        "division by zero",
			expression:  []byte(`{"expression":"2/(1 - 1)"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrDivisionByZero.Error()},
			status:      http.StatusUnprocessableEntity,
		},
		{
			name:        "bracket is not closed",
			expression:  []byte(`{"expression":"(1 + 2"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrNotClosedBracket.Error()},
			status:      http.StatusUnprocessableEntity,
		},
		{
			name:        "merger operators",
			expression:  []byte(`{"expression":"1 +* 2"}`),
			expectedErr: ErrorResponse{Error: calculator.ErrMergedOperators.Error()},
			status:      http.StatusUnprocessableEntity,
		},
	}

	//проверяем каждого теста из testCasesFail
	for _, TestCase := range testCasesFail {
		request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(TestCase.expression))
		request.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		application.СalcHandler(w, request)
		res := w.Result()
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		var actualErr ErrorResponse
		err = json.Unmarshal(data, &actualErr)
		if err != nil {
			t.Fatal(err)
		}
		//проверка совпадения ошибки с ожидаемой ошибкой
		if TestCase.expectedErr != actualErr {
			t.Fatalf("Expected error: %s, but got: %s", TestCase.expectedErr, data)
		}
		//проверка совпадения статуса ответа с ожидаемым ответом
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("Expected status: %d, but got: %d", http.StatusUnprocessableEntity, res.StatusCode)
		}
	}
}
