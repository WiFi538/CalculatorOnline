package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/wifi538/CalculatorOnline/pkg/calculator"
)

// определение порта сервера
const (
	port = ":8080"
)

type Application struct {
}

func New() *Application {
	return &Application{}
}

// структура для десериализации входящего JSON-запроса
type Request struct {
	Expression string `json:"expression"` //выражение для вычисления
}

// структура для сериализации ответа с результатом вычисления
type Response struct {
	Result float64 `json:"result"` //результат вычисления
}

// структура для сериализации ответа с ошибкой
type Error struct {
	Result string `json:"error"` //текст ошибки
}

// обработка HTTP-запросов на вычисление выражения
func СalcHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %v %v", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		//возвращаем статус 405 и соообщение об ошибке
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("Code: %v, Invalid request method", http.StatusMethodNotAllowed)
		e := Error{Result: "invalid request method"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(e)
		return
	}

	//десериализуем тело запроса в Request
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	log.Printf("Expression: %v", req)
	if err != nil {
		//если десериализация не удалась, возвращаем статус 422 и сообщение об ошибке
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("Code: %v, Invalid request body", http.StatusUnprocessableEntity)
		e := Error{Result: "invalid request body"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(e)
		return
	}

	fmt.Println(req.Expression)
	//вызываем Calc для вычисления выражения
	result, err := calculator.Calc(req.Expression)
	if err != nil {
		//если возникла ошибка, возвращаем статус 422 и сообщение об ошбке
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Printf("Code: %v, Error: %v", http.StatusUnprocessableEntity, err)
		e := Error{Result: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(e)
		return
	}

	//логируем успешный результат и возвращаем его
	log.Printf("Code: %v, Result: %v", http.StatusOK, result)
	resp := Response{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// запуск сервера на указанном порту
func (a *Application) RunServer() {
	http.HandleFunc("/api/v1/calculate", СalcHandler)

	log.Printf("Starting server on %v", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
