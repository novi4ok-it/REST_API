package handlers

import (
	"RestAPI/internal/service"
	"RestAPI/pkg/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// RegisterRequest represents user registration data
// swagger:model
type RegisterRequest struct {
	// Username for registration
	// required: true
	// example: john_doe
	Username string `json:"username" validate:"required"`

	// Password for registration
	// required: true
	// example: P@ssw0rd!
	Password string `json:"password" validate:"required"`
}

// LoginRequest represents user login data
// swagger:model
type LoginRequest struct {
	// Username for login
	// required: true
	// example: john_doe
	Username string `json:"username" validate:"required"`

	// Password for login
	// required: true
	// example: P@ssw0rd!
	Password string `json:"password" validate:"required"`
}

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register godoc
// @Summary User registration
// @Description Create new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 409 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
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

// Login godoc
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credentials"
// @Success 200 {object} responses.Response
// @Failure 400 {object} responses.Response
// @Failure 401 {object} responses.Response
// @Failure 500 {object} responses.Response
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
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
