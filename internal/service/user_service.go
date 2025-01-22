package service

import (
	"RestAPI/internal/models"
	"RestAPI/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

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
