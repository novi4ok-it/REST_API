package routes

import (
	_ "RestAPI/docs" // Импорт сгенерированной документации
	handlers "RestAPI/internal/handlers"
	repository "RestAPI/internal/repository"
	service "RestAPI/internal/service"
	"RestAPI/pkg/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
	"time"
)

// @title TodoList Management API
// @version 1.0.0
// @description Comprehensive API for managing todo lists and tasks with JWT authentication

// @contact.name API Support Team
// @contact.url https://support.todolist.com
// @contact.email support@todolist.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host api.todolist.com
// @BasePath /v1
// @schemes https http

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

// @tag.name Authentication
// @tag.description User registration and login operations

// @tag.name TodoLists
// @tag.description Operations with todo lists (create, read, update, delete)

// @tag.name Tasks
// @tag.description Operations with tasks inside todo lists (create, read, update, delete)

const (
	secretKey   = "triss-merigold"
	tokenExpiry = time.Hour * 24
)

// SetupRoutes initializes all API endpoints and middleware
// @Summary Initialize application routes
// @Description Configures all HTTP endpoints with proper middleware and security requirements
// @Tags Configuration
// @Produce json
// @Success 200 {object} responses.Response
func SetupRoutes(e *echo.Echo, db *gorm.DB) *echo.Echo {
	// Инициализация репозиториев
	todoListRepo := repository.NewTodoListRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервисов
	todoListService := service.NewTodoListService(todoListRepo)
	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo, secretKey, tokenExpiry)

	// Инициализация обработчиков
	todoListHandler := handlers.NewTodoListHandler(todoListService)
	taskHandler := handlers.NewTaskHandler(taskService, todoListService)
	authHandler := handlers.NewAuthHandler(userService)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Public routes (no authentication required)
	// Группа: Authentication
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// Protected routes (JWT authentication required)
	protected := e.Group("")
	protected.Use(middleware.JWTMiddleware(secretKey))

	// Группа: TodoLists
	protected.GET("/todolists", todoListHandler.GetTodoListHandler)
	protected.POST("/todolists", todoListHandler.PostTodoListHandler)
	protected.PATCH("/todolists/:id", todoListHandler.PatchTodoListHandler)
	protected.DELETE("/todolists/:id", todoListHandler.DeleteTodoListHandler)

	// Группа: Tasks
	protected.GET("/todolists/:list_id/tasks", taskHandler.GetTasksByListHandler)
	protected.POST("/todolists/:list_id/tasks", taskHandler.PostTaskHandler)
	protected.PATCH("/todolists/:list_id/tasks/:id", taskHandler.PatchTaskHandler)
	protected.DELETE("/todolists/:list_id/tasks/:id", taskHandler.DeleteTaskHandler)

	return e
}
