package application

import (
	"RestAPI/routes"
	"RestAPI/server"
	"RestAPI/validator"
	"context"
	"time"
)

type Application struct {
	Server *server.Server
}

func NewApplication() *Application {
	e := routes.SetupRoutes()
	e.Validator = validator.NewValidator()

	srv := server.NewServer(e, ":8080")
	return &Application{Server: srv}
}

func (a *Application) Run() {
	a.Server.Start()

	a.Server.WaitForShutdownSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.Server.Shutdown(ctx)
}
