package main

import "RestAPI/application"

func main() {
	application := application.NewApp()
	application.Run()
}
