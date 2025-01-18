package handlers

import (
	"RestAPI/config"
	"RestAPI/models"
	"RestAPI/utils"
	"errors"
	"github.com/labstack/echo"
	"gorm.io/gorm"
	"net/http"
)

// GetTodoListHandler - обработчик получения всех списков
func GetTodoListHandler(c echo.Context) error {
	var todoLists []models.TodoList

	// Использование Preload для загрузки связанных задач
	if err := config.DB.Preload("Tasks").Find(&todoLists).Error; err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not find the lists")
	}

	return c.JSON(http.StatusOK, todoLists)
}

// PostTodoListHandler - обработчик создания списка
func PostTodoListHandler(c echo.Context) error {
	// 1. Структура запроса
	type CreateTodoListRequest struct {
		Title  string `json:"title"`
		UserID int    `json:"user_id"`
	}

	// 2. Создаем экземпляр структуры
	var createListRequest CreateTodoListRequest

	// 3. Привязка JSON с полями структуры
	if err := c.Bind(&createListRequest); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not bind the request body")
	}

	// 4. Проверка user_id
	if createListRequest.UserID <= 0 {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "User ID is not valid")
	}

	// 5. Создание TodoList в БД
	newList := models.TodoList{
		Title:  createListRequest.Title,
		UserID: createListRequest.UserID,
	}

	if err := config.DB.Create(&newList).Error; err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not create the list")
	}

	// 6. Возвращаем успешный ответ
	return utils.JSONResponse(c, http.StatusCreated, "ok", "TodoList was successfully created")
}

// PatchTodoListHandler - обработчик обновления списка
func PatchTodoListHandler(c echo.Context) error {
	id, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}

	// 1. Структура для запроса
	type UpdateTodoListRequest struct {
		Title string `json:"title"`
	}

	// 2. Создаем экземпляр
	var updateRequest UpdateTodoListRequest

	// 3. Привязка данных из запроса
	if err := c.Bind(&updateRequest); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not bind the request body")
	}

	// 4. Проверка title
	if updateRequest.Title == "" {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Title is required")
	}

	// 5. Обновляем данные в базе данных
	var todoList models.TodoList
	if err := config.DB.First(&todoList, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusBadRequest, "error", "TodoList with this ID does not exist")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Error finding TodoList")
	}

	todoList.Title = updateRequest.Title
	if err := config.DB.Save(&todoList).Error; err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not update the list")
	}

	return utils.JSONResponse(c, http.StatusOK, "ok", "TodoList was successfully updated")
}

func DeleteTodoListHandler(c echo.Context) error {
	id, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}

	var todoList models.TodoList
	if err := config.DB.Preload("Tasks").First(&todoList, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusBadRequest, "error", "TodoList with this ID does not exist")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Error finding TodoList")
	}
	// 1. Удаление связанных тасков
	if len(todoList.Tasks) > 0 {
		if err := config.DB.Delete(&todoList.Tasks).Error; err != nil {
			return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not delete tasks of this list")
		}
	}
	// 2. Удаление списка Todo
	if err := config.DB.Delete(&todoList).Error; err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not delete the list")
	}

	return utils.JSONResponse(c, http.StatusOK, "ok", "List was successfully deleted")
}

func GetTasksByListHandler(c echo.Context) error {
	listID, err := utils.GetParam(c, "list_id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad list ID")
	}

	var tasks []models.Task

	if err := config.DB.Where("list_id = ?", listID).Find(&tasks).Error; err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not find the tasks")
	}

	return c.JSON(http.StatusOK, tasks)
}

// PostTaskHandler - обработчик создания таска
func PostTaskHandler(c echo.Context) error {
	// 1. Получение list_id из URL
	listID, err := utils.GetParam(c, "list_id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad list ID")
	}

	// 2. Проверка, что лист с таким ID существует
	var list models.TodoList
	if err := config.DB.First(&list, listID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusBadRequest, "error", "TodoList with this ID does not exist")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Error checking TodoList")
	}

	// 3. Структура для запроса
	type CreateTaskRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var taskReq CreateTaskRequest
	if err := c.Bind(&taskReq); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not bind the request body")
	}

	// 4. Создание таска
	newTask := models.Task{
		Title:       taskReq.Title,
		Description: taskReq.Description,
		ListID:      listID,
		Completed:   false,
	}
	// 5. Создаем запись в базе
	if err := config.DB.Create(&newTask).Error; err != nil {
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Could not create the task")
	}
	return utils.JSONResponse(c, http.StatusCreated, "ok", "Task was successfully created")
}

// PatchTaskHandler - обработчик обновления таска
func PatchTaskHandler(c echo.Context) error {
	// 1. Получаем task id из URL
	taskID, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad task ID")
	}

	// 2. Структура для запроса
	type UpdateTaskRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	// 3. Создаем экземпляр структуры
	var updateRequest UpdateTaskRequest

	// 4. Связывание данных из запроса
	if err := c.Bind(&updateRequest); err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not bind the request body")
	}

	// 5. Обновляем данные
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusBadRequest, "error", "Task with this ID does not exist")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Error finding Task")
	}

	if updateRequest.Title != "" {
		task.Title = updateRequest.Title
	}
	if updateRequest.Description != "" {
		task.Description = updateRequest.Description
	}

	if c.Request().Body != nil {
		task.Completed = updateRequest.Completed
	}

	if err := config.DB.Save(&task).Error; err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not update the task")
	}

	return utils.JSONResponse(c, http.StatusOK, "ok", "Task was successfully updated")
}

// DeleteTaskHandler - обработчик удаления таска
func DeleteTaskHandler(c echo.Context) error {
	id, err := utils.GetParam(c, "id")
	if err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Bad ID")
	}

	var task models.Task
	if err := config.DB.First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.JSONResponse(c, http.StatusBadRequest, "error", "Task with this ID does not exist")
		}
		return utils.JSONResponse(c, http.StatusInternalServerError, "error", "Error finding Task")
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		return utils.JSONResponse(c, http.StatusBadRequest, "error", "Could not delete the task")
	}
	return utils.JSONResponse(c, http.StatusOK, "ok", "Task was successfully deleted")
}
