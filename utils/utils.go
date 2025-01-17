package utils

import (
	"RestAPI/responses"
	"github.com/labstack/echo"
	"strconv"
)

func JSONResponse(c echo.Context, statusCode int, status, message string) error {
	return c.JSON(statusCode, responses.Response{
		Status:  status,
		Message: message,
	})
}

func GetIDParam(c echo.Context) (int, error) {
	idParam := c.Param("id")
	return strconv.Atoi(idParam)
}
