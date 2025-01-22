package service

import (
	"RestAPI/internal/models"
	"RestAPI/internal/repository"
	"errors"
)

type TaskService interface {
	GetAllTasksForList(listID int, userID int) ([]models.Task, error)
	GetTaskByID(taskID int, userID int) (*models.Task, error)
	CreateTask(title string, description string, listID int) error
	UpdateTask(taskID int, userID int, title string, description string, isCompleted *bool) error
	DeleteTask(taskID int, userID int) error
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

type taskService struct {
	repo repository.TaskRepository
}

func (s *taskService) GetAllTasksForList(listID int, userID int) ([]models.Task, error) {
	if listID <= 0 {
		return nil, errors.New("invalid list ID")
	}
	return s.repo.GetAllTasksForThisList(listID, userID)
}

func (s *taskService) GetTaskByID(taskID int, userID int) (*models.Task, error) {
	if taskID <= 0 {
		return nil, errors.New("invalid task ID")
	}
	return s.repo.GetTaskByID(taskID, userID)
}

func (s *taskService) CreateTask(title string, description string, listID int) error {
	if listID <= 0 {
		return errors.New("invalid list ID")
	}
	if title == "" {
		return errors.New("task title cannot be empty")
	}
	if description == "" {
		return errors.New("task description cannot be empty")
	}

	task := &models.Task{
		Title:       title,
		Description: description,
		ListID:      listID,
		Completed:   false,
	}

	err := s.repo.CreateTask(task)
	return err
}

func (s *taskService) UpdateTask(taskID int, userID int, title string, description string, isCompleted *bool) error {
	task, err := s.GetTaskByID(taskID, userID)
	if err != nil {
		return err
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}

	if isCompleted != nil {
		task.Completed = *isCompleted
	}

	return s.repo.UpdateTask(task)
}

func (s *taskService) DeleteTask(taskID int, userID int) error {
	task, err := s.GetTaskByID(taskID, userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteTask(task)
}
