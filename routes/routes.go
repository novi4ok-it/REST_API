package routes

import (
	"RestAPI/handlers"
	"github.com/labstack/echo"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/messages", handlers.GetHandler)
	e.POST("/messages", handlers.PostHandler)
	e.PATCH("/messages/:id", handlers.PatchHandler)
	e.DELETE("/messages/:id", handlers.DeleteHandler)
}
