package handlers

import (
	"RestAPI/service"
	"RestAPI/utils"
	"errors"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

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
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	tasks, err := h.taskService.GetAllTasksForList(listID, int(userID))
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
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	_, err = h.todoListService.GetListByID(listID, int(userID)) //На этом уровне идёт проверка, принадлежит ли данный лист этому пользователю
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
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	err = h.taskService.UpdateTask(taskID, int(userID), req.Title, req.Description, req.Completed)
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
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	err = h.taskService.DeleteTask(taskID, int(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusNotFound, "error", "Task with this ID does not exist")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not delete the task")
	}

	return utils.JSONResponse(c, http.StatusOK, "ok", "Task deleted successfully")
}
