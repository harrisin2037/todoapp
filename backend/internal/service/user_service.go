package service

import (
	"errors"

	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/repository"
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

	token, err := models.GenerateToken(user)
	if err != nil {
		return nil, "", ErrFailedToGenerateJWT
	}

	return user, token, nil
}

func (s *UserService) GetUsersByIDs(ids []uint) ([]models.User, error) {
	return s.repo.GetUsersByIDs(ids)
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

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateUser(userID uint, username, email, password, role string) error {

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if username != "" {
		user.Username = username
	}

	if email != "" {
		user.Email = email
	}

	if password != "" {
		err = user.SetPassword(password)
		if err != nil {
			return err
		}
	}

	if role != "" {
		newRole, err := models.ParseRole(role)
		if err != nil {
			return err
		}
		user.Role = newRole
	}

	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(userID uint) error {
	_, err := s.repo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(userID)
}

func (s *UserService) CreateUser(username, email, password string, role models.Role) error {
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
		Role:     role,
	}
	err := newUser.SetPassword(password)
	if err != nil {
		return err
	}

	return s.repo.Create(newUser)
}
