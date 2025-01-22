package repository

import (
	"RestAPI/internal/models"
	"gorm.io/gorm"
)

type TodoListRepository interface {
	GetAllLists(userID int) ([]models.TodoList, error)
	GetListByID(listID int, userID int) (*models.TodoList, error)
	CreateList(todoList *models.TodoList) error
	UpdateList(todoList *models.TodoList) error
	DeleteList(todoList *models.TodoList) error
	DeleteAllTasksForThisList(todoList *models.TodoList) error
}

type todoListRepository struct {
	DB *gorm.DB
}

func NewTodoListRepository(db *gorm.DB) TodoListRepository {
	return &todoListRepository{DB: db}
}

func (r *todoListRepository) GetAllLists(userID int) ([]models.TodoList, error) {
	var todoLists []models.TodoList
	err := r.DB.Preload("Tasks").Where("user_id = ?", userID).Find(&todoLists).Error
	return todoLists, err
}

func (r *todoListRepository) GetListByID(listID int, userID int) (*models.TodoList, error) {
	var todoList models.TodoList
	err := r.DB.Preload("Tasks").Where("user_id = ?", userID).First(&todoList, listID).Error
	return &todoList, err
}

func (r *todoListRepository) CreateList(todoList *models.TodoList) error {
	return r.DB.Create(&todoList).Error
}

func (r *todoListRepository) UpdateList(todoList *models.TodoList) error {
	return r.DB.Save(&todoList).Error
}

func (r *todoListRepository) DeleteList(todoList *models.TodoList) error {
	return r.DB.Delete(&todoList).Error
}

func (r *todoListRepository) DeleteAllTasksForThisList(todoList *models.TodoList) error {
	return r.DB.Delete(&todoList.Tasks).Error
}
