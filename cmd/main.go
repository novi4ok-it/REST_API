package main

import (
	"RestAPI/internal/application"
)

func main() {
	application := application.NewApp()
	application.Run()
}
