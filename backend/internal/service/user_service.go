package service

import (
	"errors"

	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/repository"
	"github.com/harrisin2037/todoapp/internal/utils"
)

type UserService struct {
	repo *repository.UserRepository
}

var (
	ErrUserAlreadyExists   = errors.New("username already exists")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrFailedToGenerateJWT = errors.New("failed to generate JWT token")
	ErrSameRole            = errors.New("user already has this role")
)

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, email, password string) error {
	existingUser, _ := s.repo.FindByUsername(username)
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	existingUser, _ = s.repo.FindByEmail(email)
	if existingUser != nil {
		return ErrEmailAlreadyExists
	}

	newUser := &models.User{
		Username: username,
		Email:    email,
		Role:     models.RoleUser,
	}
	err := newUser.SetPassword(password)
	if err != nil {
		return err
	}

	return s.repo.Create(newUser)
}

func (s *UserService) Login(username, password string) (*models.User, string, error) {

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if !user.CheckPassword(password) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		return nil, "", ErrFailedToGenerateJWT
	}

	return user, token, nil
}

func (s *UserService) ChangeRole(userID uint, newRole models.Role) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.Role == newRole {
		return ErrSameRole
	}

	return s.repo.UpdateRole(userID, newRole)
}
