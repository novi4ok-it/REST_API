package main

import (
	"RestAPI/routes"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"log"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	e := echo.New()
	e = routes.SetupRoutes()
	e.Validator = &CustomValidator{validator: validator.New()}
	// Запускаем сервер
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
