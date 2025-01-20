package handlers

import (
	"RestAPI/service"
	"RestAPI/utils"
	"fmt"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

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
