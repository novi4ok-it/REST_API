package main

import (
	"RestAPI/config"
	"RestAPI/routes"
	"github.com/labstack/echo"
)

func main() {
	config.InitDB()
	e := echo.New()
	routes.RegisterRoutes(e)
	e.Start(":8080")
}
