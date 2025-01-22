package routes

import (
	handlers2 "RestAPI/internal/handlers"
	repository2 "RestAPI/internal/repository"
	service2 "RestAPI/internal/service"
	"RestAPI/pkg/middleware"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"time"
)

const secretKey = "triss-merigold"
const tokenExpiry = time.Hour * 24

func SetupRoutes(e *echo.Echo, db *gorm.DB) *echo.Echo {

	todoListRepo := repository2.NewTodoListRepository(db)
	taskRepo := repository2.NewTaskRepository(db)
	userRepo := repository2.NewUserRepository(db)

	todoListService := service2.NewTodoListService(todoListRepo)
	taskService := service2.NewTaskService(taskRepo)
	userService := service2.NewUserService(userRepo, secretKey, tokenExpiry)

	todoListHandler := handlers2.NewTodoListHandler(todoListService)
	taskHandler := handlers2.NewTaskHandler(taskService, todoListService)
	authHandler := handlers2.NewAuthHandler(userService)

	// Маршрут для регистрации
	e.POST("/register", authHandler.Register)

	// Маршрут для логина
	e.POST("/login", authHandler.Login)

	// Защищенные маршруты
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware(secretKey))

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
