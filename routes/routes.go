package routes

import (
	"RestAPI/config"
	"RestAPI/handlers"
	"RestAPI/repository"
	"RestAPI/service"
	"github.com/labstack/echo"
)

func SetupRoutes() *echo.Echo {
	e := echo.New()

	config.InitDB()
	db := config.DB

	todoListRepo := repository.NewTodoListRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	todoListService := service.NewTodoListService(todoListRepo)
	taskService := service.NewTaskService(taskRepo)

	todoListHandler := handlers.NewTodoListHandler(todoListService)
	taskHandler := handlers.NewTaskHandler(taskService, todoListService)

	e.GET("/todolists", todoListHandler.GetTodoListHandler)
	e.POST("/todolists", todoListHandler.PostTodoListHandler)
	e.PATCH("/todolists/:id", todoListHandler.PatchTodoListHandler)
	e.DELETE("/todolists/:id", todoListHandler.DeleteTodoListHandler)

	e.GET("/todolists/:list_id/tasks", taskHandler.GetTasksByListHandler)
	e.POST("/todolists/:list_id/tasks", taskHandler.PostTaskHandler)
	e.PATCH("/todolists/:list_id/tasks/:id", taskHandler.PatchTaskHandler)
	e.DELETE("/todolists/:list_id/tasks/:id", taskHandler.DeleteTaskHandler)

	return e
}
