package service

import (
	"RestAPI/models"
	"RestAPI/repository"
	"errors"
)

type TodoListService interface {
	GetAllLists(userID int) ([]models.TodoList, error)
	GetListByID(listID int, userID int) (*models.TodoList, error)
	CreateList(title string, userID int) error
	UpdateList(listID int, userID int, title string) error
	DeleteList(listID int, userID int) error
}

type todoListService struct {
	repo repository.TodoListRepository
}

func NewTodoListService(repo repository.TodoListRepository) TodoListService {
	return &todoListService{repo: repo}
}

func (s *todoListService) GetAllLists(userID int) ([]models.TodoList, error) {
	return s.repo.GetAllLists(userID)
}

func (s *todoListService) GetListByID(listID int, userID int) (*models.TodoList, error) {
	if listID <= 0 {
		return nil, errors.New("invalid list ID")
	}
	return s.repo.GetListByID(listID, userID)
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

func (s *todoListService) UpdateList(listID int, userID int, title string) error {
	list, err := s.GetListByID(listID, userID)
	if err != nil {
		return err
	}
	list.Title = title
	return s.repo.UpdateList(list)
}

func (s *todoListService) DeleteList(listID int, userID int) error {
	list, err := s.GetListByID(listID, userID)
	if err != nil {
		return err
	}
	if len(list.Tasks) > 0 {
		s.repo.DeleteAllTasksForThisList(list)
	}
	return s.repo.DeleteList(list)
}
