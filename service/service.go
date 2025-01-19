package service

import (
	"RestAPI/models"
	"RestAPI/repository"
	"errors"
)

type TodoListService interface {
	GetAllLists() ([]models.TodoList, error)
	GetListByID(id int) (*models.TodoList, error)
	CreateList(title string, userID int) error
	UpdateList(id int, title string) error
	DeleteList(id int) error
}

type todoListService struct {
	repo repository.TodoListRepository
}

func NewTodoListService(repo repository.TodoListRepository) TodoListService {
	return &todoListService{repo: repo}
}

func (s *todoListService) GetAllLists() ([]models.TodoList, error) {
	return s.repo.GetAllLists()
}

func (s *todoListService) GetListByID(id int) (*models.TodoList, error) {
	return s.repo.GetListByID(id)
}

func (s *todoListService) CreateList(title string, userID int) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}
	list := &models.TodoList{
		Title:  title,
		UserID: userID,
	}
	err := s.repo.CreateList(list)
	return err
}

func (s *todoListService) UpdateList(id int, title string) error {
	list, err := s.repo.GetListByID(id)
	if err != nil {
		return err
	}
	list.Title = title
	return s.repo.UpdateList(list)
}

func (s *todoListService) DeleteList(id int) error {
	list, err := s.repo.GetListByID(id)
	if err != nil {
		return err
	}
	if len(list.Tasks) > 0 {
		s.repo.DeleteAllTasksForThisList(list)
	}
	return s.repo.DeleteList(list)
}

type TaskService interface {
	GetAllTasksForList(listID int) ([]models.Task, error)
	GetTaskByID(taskID int) (*models.Task, error)
	CreateTask(title string, description string, listID int) error
	UpdateTask(taskID int, title string, description string, isCompleted *bool) error
	DeleteTask(taskID int) error
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

type taskService struct {
	repo repository.TaskRepository
}

func (s *taskService) GetAllTasksForList(listID int) ([]models.Task, error) {
	if listID <= 0 {
		return nil, errors.New("invalid list ID")
	}
	return s.repo.GetAllTasksForThisList(listID)
}

func (s *taskService) GetTaskByID(taskID int) (*models.Task, error) {
	if taskID <= 0 {
		return nil, errors.New("invalid task ID")
	}
	return s.repo.GetTaskByID(taskID)
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

func (s *taskService) UpdateTask(taskID int, title string, description string, isCompleted *bool) error {
	task, err := s.repo.GetTaskByID(taskID)
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

func (s *taskService) DeleteTask(taskID int) error {
	task, err := s.repo.GetTaskByID(taskID)
	if err != nil {
		return err
	}

	return s.repo.DeleteTask(task)
}
