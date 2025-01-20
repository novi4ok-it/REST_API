package handlers

import (
	"RestAPI/service"
	"RestAPI/utils"
	"errors"
	"github.com/labstack/echo"
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

func (h *todoListHandler) PostTodoListHandler(c echo.Context) error {
	type CreateTodoListRequest struct {
		Title string `json:"title"`
	}

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

func (h *todoListHandler) PatchTodoListHandler(c echo.Context) error {
	listID, err := utils.GetParam(c, "id")
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
