package handlers

import (
	"RestAPI/service"
	"RestAPI/utils"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type TodoListHandler interface {
	GetTodoListHandler(c echo.Context) error
	PostTodoListHandler(c echo.Context) error
	PatchTodoListHandler(c echo.Context) error
	DeleteTodoListHandler(c echo.Context) error
}

type todoListHandler struct {
	todoListService service.TodoListService
}

func NewTodoListHandler(todoListService service.TodoListService) TodoListHandler {
	return &todoListHandler{todoListService: todoListService}
}

func (h *todoListHandler) GetTodoListHandler(c echo.Context) error {
	todoLists, err := h.todoListService.GetAllLists()
	if err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not fetch todo lists")
	}
	return c.JSON(http.StatusOK, todoLists)
}

func (h *todoListHandler) PostTodoListHandler(c echo.Context) error {
	type CreateTodoListRequest struct {
		Title  string `json:"title"`
		UserID int    `json:"user_id"`
	}

	var req CreateTodoListRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid request body")
	}

	err := h.todoListService.CreateList(req.Title, req.UserID)
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", err.Error())
	}
	return utils.JSONResponse(c, http.StatusCreated, "ok", "TodoList was successfully created")
}

func (h *todoListHandler) PatchTodoListHandler(c echo.Context) error {
	id, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}

	type UpdateTodoListRequest struct {
		Title string `json:"title"`
	}

	var req UpdateTodoListRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid request body")
	}

	if err := h.todoListService.UpdateList(id, req.Title); err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "List updated successfully")
}

func (h *todoListHandler) DeleteTodoListHandler(c echo.Context) error {
	id, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}

	if err := h.todoListService.DeleteList(id); err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "List deleted successfully")
}

type TaskHandler interface {
	GetTasksByListHandler(c echo.Context) error
	PostTaskHandler(c echo.Context) error
	PatchTaskHandler(c echo.Context) error
	DeleteTaskHandler(c echo.Context) error
}

type taskHandler struct {
	taskService     service.TaskService
	todoListService service.TodoListService
}

func NewTaskHandler(taskService service.TaskService, todoListService service.TodoListService) TaskHandler {
	return &taskHandler{taskService: taskService,
		todoListService: todoListService}
}

func (h *taskHandler) GetTasksByListHandler(c echo.Context) error {
	listID, err := utils.GetParam(c, "list_id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid list ID")
	}

	tasks, err := h.taskService.GetAllTasksForList(listID)
	if err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not fetch tasks for the list")
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *taskHandler) PostTaskHandler(c echo.Context) error {
	listID, err := utils.GetParam(c, "list_id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid list ID")
	}

	_, err = h.todoListService.GetListByID(listID)
	if err != nil {
		return utils.JSONResponse(c, http.StatusNotFound, "error", "TodoList with this ID does not exist")
	}

	type CreateTaskRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var req CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid request body")
	}

	err = h.taskService.CreateTask(req.Title, req.Description, listID)
	if err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not create task")
	}

	return utils.JSONResponse(c, http.StatusCreated, "ok", "Task was successfully created")
}

func (h *taskHandler) PatchTaskHandler(c echo.Context) error {
	taskID, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid task ID")
	}

	type UpdateTaskRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   *bool  `json:"completed"`
	}

	var req UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid request body")
	}

	err = h.taskService.UpdateTask(taskID, req.Title, req.Description, req.Completed)
	if err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "ok", "Task updated successfully")
}

func (h *taskHandler) DeleteTaskHandler(c echo.Context) error {
	taskID, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid task ID")
	}

	err = h.taskService.DeleteTask(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusNotFound, "error", "Task with this ID does not exist")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not delete the task")
	}

	return utils.JSONResponse(c, http.StatusOK, "ok", "Task deleted successfully")
}

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(c echo.Context) error {

	type RegisterRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid input")
	}

	if err := c.Validate(&req); err != nil {
		log.Println("Validation error:", err) // Добавьте логирование ошибки
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid input")
	}

	err := h.userService.RegisterUser(req.Username, req.Password)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			return utils.JSONResponse(c, http.StatusConflict, "error", "Username already exists")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Failed to register user")
	}

	return utils.JSONResponse(c, http.StatusCreated, "ok", "User registered successfully")
}

func (h *AuthHandler) Login(c echo.Context) error {
	type LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid input")
	}

	if err := c.Validate(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid input")
	}

	token, err := h.userService.LoginUser(req.Username, req.Password)
	if err != nil {
		if err == service.ErrUserNotFound || err == service.ErrInvalidCredentials {
			return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid credentials")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Failed to login")
	}

	return utils.JSONResponse(c, http.StatusOK, "ok", fmt.Sprintf("%s", token))
}
