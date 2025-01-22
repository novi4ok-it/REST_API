package handlers

import (
	service "RestAPI/internal/service"
	"RestAPI/pkg/utils"
	"errors"
	"github.com/labstack/echo/v4"
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

// CreateTaskRequest represents data for creating new task
// swagger:model
type CreateTaskRequest struct {
	// Title of the task
	// required: true
	// example: Buy groceries
	Title string `json:"title"`

	// Description of the task
	// example: Milk, eggs, bread
	Description string `json:"description"`
}

// UpdateTaskRequest model
// swagger:model
type UpdateTaskRequest struct {
	// New title for the task
	// example: Buy organic milk
	Title string `json:"title"`

	// New description for the task
	// example: 2 liters of organic milk
	Description string `json:"description"`

	// New completion status
	// example: true
	Completed *bool `json:"completed"`
}

// GetTasksByListHandler godoc
// @Summary Get tasks by list
// @Description Get all tasks for specified todo list
// @Tags tasks
// @Security Bearer
// @Produce json
// @Param list_id path int true "Todo List ID"
// @Success 200 {array} models.Task
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists/{list_id}/tasks [get]
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

// PostTaskHandler godoc
// @Summary Create new task
// @Description Create new task in specified todo list
// @Tags tasks
// @Security Bearer
// @Accept json
// @Produce json
// @Param list_id path int true "Todo List ID"
// @Param request body handlers.CreateTaskRequest true "Task data"
// @Success 201 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists/{list_id}/tasks [post]
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

// PatchTaskHandler godoc
// @Summary Update task
// @Description Update task details (title, description, completed status)
// @Tags tasks
// @Security Bearer
// @Accept json
// @Produce json
// @Param list_id path int true "Todo List ID"
// @Param id path int true "Task ID"
// @Param request body handlers.UpdateTaskRequest true "Task update data"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists/{list_id}/tasks/{id} [patch]
func (h *taskHandler) PatchTaskHandler(c echo.Context) error {
	taskID, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid task ID")
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

// DeleteTaskHandler godoc
// @Summary Delete task
// @Description Delete task from todo list
// @Tags tasks
// @Security Bearer
// @Produce json
// @Param list_id path int true "Todo List ID"
// @Param id path int true "Task ID"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists/{list_id}/tasks/{id} [delete]
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
