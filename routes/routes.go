package routes

import (
	"RestAPI/config"
	"RestAPI/handlers"
	"RestAPI/middleware"
	"RestAPI/repository"
	"RestAPI/service"
	"github.com/labstack/echo"
	"time"
)

func SetupRoutes() *echo.Echo {
	e := echo.New()

	config.InitDB()
	db := config.DB

	todoListRepo := repository.NewTodoListRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	todoListService := service.NewTodoListService(todoListRepo)
	taskService := service.NewTaskService(taskRepo)
	const secretKey = "triss-merigold"
	const tokenExpiry = time.Hour * 24
	userService := service.NewUserService(userRepo, secretKey, tokenExpiry)

	todoListHandler := handlers.NewTodoListHandler(todoListService)
	taskHandler := handlers.NewTaskHandler(taskService, todoListService)
	authHandler := handlers.NewAuthHandler(userService)

	// Маршрут для регистрации
	e.POST("/register", authHandler.Register)

	// Маршрут для логина
	e.POST("/login", authHandler.Login)

	// Защищенные маршруты
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware(secretKey)) // Применение middleware для проверки JWT

	// Теперь эти маршруты защищены и требуют токена
	protected.GET("/todolists", todoListHandler.GetTodoListHandler)
	protected.POST("/todolists", todoListHandler.PostTodoListHandler)
	protected.PATCH("/todolists/:id", todoListHandler.PatchTodoListHandler)
	protected.DELETE("/todolists/:id", todoListHandler.DeleteTodoListHandler)

	protected.GET("/todolists/:list_id/tasks", taskHandler.GetTasksByListHandler)
	protected.POST("/todolists/:list_id/tasks", taskHandler.PostTaskHandler)
	protected.PATCH("/todolists/:list_id/tasks/:id", taskHandler.PatchTaskHandler)
	protected.DELETE("/todolists/:list_id/tasks/:id", taskHandler.DeleteTaskHandler)

	return e
}
