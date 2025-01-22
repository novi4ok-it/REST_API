package utils

import (
	"RestAPI/internal/responses"
	"errors"
	"github.com/labstack/echo"
	"strconv"
)

func JSONResponse(c echo.Context, statusCode int, status, message string) error {
	return c.JSON(statusCode, responses.Response{
		Status:  status,
		Message: message,
	})
}

func GetParam(c echo.Context, key string) (int, error) {
	idParam := c.Param(key)
	if idParam == "" {
		return 0, errors.New("id parameter is missing")
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
