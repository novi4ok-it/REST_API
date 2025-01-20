package repository

import (
	"RestAPI/models"
	"gorm.io/gorm"
)

type TaskRepository interface {
	GetAllTasksForThisList(listID int, userID int) ([]models.Task, error)
	GetTaskByID(taskID int, userID int) (*models.Task, error)
	CreateTask(task *models.Task) error
	UpdateTask(task *models.Task) error
	DeleteTask(task *models.Task) error
}

type taskRepository struct {
	DB *gorm.DB
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
