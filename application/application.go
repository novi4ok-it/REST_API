package application

import (
	"RestAPI/config"
	"RestAPI/routes"
	"RestAPI/server"
	"RestAPI/validator"
	"context"
	"github.com/labstack/echo"
	"time"
)

type Application struct {
	Server *server.Server
}

func NewApp() *Application {
	e := echo.New()
	config.LoadConfig()
	config.InitDB()
	db := config.DB
	routes.SetupRoutes(e, db)
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
