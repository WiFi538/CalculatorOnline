package main

import (
	"github.com/wifi538/CalculatorOnline/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
