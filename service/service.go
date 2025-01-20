package service

import (
	"RestAPI/models"
	"RestAPI/repository"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
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

type UserService interface {
	RegisterUser(username, password string) error
	LoginUser(username, password string) (string, error)
}

type userService struct {
	repo        repository.UserRepository
	secretKey   string
	tokenExpiry time.Duration
}

func NewUserService(repo repository.UserRepository, secretKey string, tokenExpiry time.Duration) UserService {
	return &userService{repo: repo, secretKey: secretKey, tokenExpiry: tokenExpiry}
}

func (s *userService) RegisterUser(username, password string) error {
	existingUser, err := s.repo.FindByUsername(username)
	if err == nil && existingUser != nil {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrHashingPassword
	}

	newUser := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return s.repo.CreateUser(newUser)
}

func (s *userService) LoginUser(username, password string) (string, error) {
	// Находим пользователя в базе
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return "", ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(), // Используем константу
	})

	tokenString, err := token.SignedString([]byte(s.secretKey)) // Используем константу
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrHashingPassword    = errors.New("failed to hash password")
)
