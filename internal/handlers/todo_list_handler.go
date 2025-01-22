package handlers

import (
	"RestAPI/internal/service"
	"RestAPI/pkg/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

type CreateTodoListRequest struct {
	// Title of the todo list
	// required: true
	// example: My Shopping List
	Title string `json:"title"`
}

// UpdateTodoListRequest model
// swagger:model
type UpdateTodoListRequest struct {
	// New title for the list
	// required: true
	// example: Updated Shopping List
	Title string `json:"title"`
}

// GetTodoListHandler godoc
// @Summary Get all todo lists
// @Description Retrieve all todo lists for authenticated user
// @Tags todolists
// @Security Bearer
// @Produce json
// @Success 200 {array} models.TodoList
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists [get]
func (h *todoListHandler) GetTodoListHandler(c echo.Context) error {
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	todoLists, err := h.todoListService.GetAllLists(int(userID))
	if err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not fetch todo lists")
	}
	return c.JSON(http.StatusOK, todoLists)
}

// PostTodoListHandler godoc
// @Summary Create new todo list
// @Description Create new todo list for authenticated user
// @Tags todolists
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body handlers.CreateTodoListRequest true "List data"
// @Success 201 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists [post]
func (h *todoListHandler) PostTodoListHandler(c echo.Context) error {
	var req CreateTodoListRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid request body")
	}
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	err := h.todoListService.CreateList(req.Title, int(userID))
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", err.Error())
	}
	return utils.JSONResponse(c, http.StatusCreated, "ok", "TodoList was successfully created")
}

// PatchTodoListHandler godoc
// @Summary Update todo list
// @Description Update title of existing todo list
// @Tags todolists
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Todo List ID"
// @Param request body handlers.UpdateTodoListRequest true "New title"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists/{id} [patch]
func (h *todoListHandler) PatchTodoListHandler(c echo.Context) error {
	listID, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}

	var req UpdateTodoListRequest
	if err := c.Bind(&req); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid request body")
	}
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	if err := h.todoListService.UpdateList(listID, int(userID), req.Title); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusNotFound, "error", "Todo list not found")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "List updated successfully")
}

// DeleteTodoListHandler godoc
// @Summary Delete todo list
// @Description Delete todo list and all its tasks
// @Tags todolists
// @Security Bearer
// @Produce json
// @Param id path int true "Todo List ID"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 404 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /todolists/{id} [delete]
func (h *todoListHandler) DeleteTodoListHandler(c echo.Context) error {
	listID, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}
	userID, ok := c.Get("user_id").(float64) // JWT Claims возвращают float64 для чисел
	if !ok {
		return utils.JSONResponse(c, http.StatusUnauthorized, "error", "Invalid or missing token")
	}
	if err := h.todoListService.DeleteList(listID, int(userID)); err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "List deleted successfully")
}
