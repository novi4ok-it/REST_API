package main

import (
	"RestAPI/config"
	"RestAPI/routes"
	"github.com/labstack/echo"
	"log"
)

func main() {
	config.InitDB()
	e := echo.New()
	routes.RegisterRoutes(e)
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
