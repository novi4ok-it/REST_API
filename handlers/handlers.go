package handlers

import (
	"RestAPI/config"
	"RestAPI/dbutils"
	"RestAPI/models"
	"RestAPI/utils"
	"github.com/labstack/echo"
	"net/http"
)

func GetHandler(c echo.Context) error {
	var messages []models.Message
	if err := dbutils.FindAll(config.DB, &messages); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not find the messages")
	}
	return c.JSON(http.StatusOK, messages)
}

func PostHandler(c echo.Context) error {
	var message models.Message
	if err := c.Bind(&message); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not add the message")
	}
	if err := dbutils.Create(config.DB, &message); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not create the message")
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "Message was successfully added")
}

func PatchHandler(c echo.Context) error {
	id, err := utils.GetIDParam(c)
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}
	var updatedMessage models.Message
	if err := c.Bind(&updatedMessage); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not update the message")
	}
	if err := dbutils.UpdateByID[models.Message](config.DB, id, "text", updatedMessage.Text); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Invalid input")
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "Message was successfully updated")
}

func DeleteHandler(c echo.Context) error {
	id, err := utils.GetIDParam(c)
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}
	if err := dbutils.DeleteByID[models.Message](config.DB, id); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not delete the message")
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "Message was successfully deleted")
}
