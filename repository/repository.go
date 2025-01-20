package repository

import (
	"RestAPI/models"
	"errors"
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

type taskRepository struct {
	DB *gorm.DB
}

type TaskRepository interface {
	GetAllTasksForThisList(listID int, userID int) ([]models.Task, error)
	GetTaskByID(taskID int, userID int) (*models.Task, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(task *models.Task) error
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{DB: db}
}

func (r *taskRepository) GetAllTasksForThisList(listID int, userID int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.DB.Joins("JOIN todo_lists ON todo_lists.id = tasks.list_id").
		Where("todo_lists.user_id = ?", userID).
		Where("tasks.list_id = ?", listID).
		Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetTaskByID(taskID int, userID int) (*models.Task, error) {
	var task models.Task
	err := r.DB.Joins("JOIN todo_lists ON todo_lists.id = tasks.list_id").
		Where("todo_lists.user_id = ?", userID).
		Where("tasks.id = ?", taskID).
		First(&task).Error
	return &task, err
}

func (r *taskRepository) CreateTask(task *models.Task) error {
	return r.DB.Create(&task).Error
}

func (r *taskRepository) UpdateTask(task *models.Task) error {
	return r.DB.Save(&task).Error
}

func (r *taskRepository) DeleteTask(task *models.Task) error {
	return r.DB.Delete(&task).Error
}

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	CreateUser(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}
