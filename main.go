package main

import (
	"RestAPI/routes"
	"log"
)

func main() {
	e := routes.SetupRoutes()
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
