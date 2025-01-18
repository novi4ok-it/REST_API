package routes

import (
	"RestAPI/handlers"
	"github.com/labstack/echo"
)

func RegisterRoutes(e *echo.Echo) {
	// TodoList маршруты
	e.GET("/todolists", handlers.GetTodoListHandler)
	e.POST("/todolists", handlers.PostTodoListHandler)
	e.PATCH("/todolists/:id", handlers.PatchTodoListHandler)
	e.DELETE("/todolists/:id", handlers.DeleteTodoListHandler)
	// Task маршруты
	e.GET("/todolists/:list_id/tasks", handlers.GetTasksByListHandler)
	e.POST("/todolists/:list_id/tasks", handlers.PostTaskHandler)
	e.PATCH("/todolists/:list_id/tasks/:id", handlers.PatchTaskHandler)
	e.DELETE("/todolists/:list_id/tasks/:id", handlers.DeleteTaskHandler)

}
