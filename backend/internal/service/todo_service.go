package service

import (
	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/repository"
)

type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(todo *models.Todo, assigneeIDs []uint) error {
	return s.repo.Create(todo, assigneeIDs)
}

func (s *TodoService) GetTodosByAdmin(statuses []string, sortBy, order string) ([]models.Todo, error) {
	return s.repo.GetListByAdmin(statuses, sortBy, order)
}

func (s *TodoService) GetTodos(statuses []string, sortBy, order string, userID uint) ([]models.Todo, error) {
	return s.repo.GetList(statuses, sortBy, order, userID)
}

func (s *TodoService) GetTodo(id uint) (*models.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *TodoService) UpdateTodo(todo *models.Todo) error {
	return s.repo.Update(todo)
}

func (s *TodoService) DeleteTodo(id uint) error {
	return s.repo.Delete(id)
}

func (s *TodoService) AssignUser(todoID, userID uint) error {
	return s.repo.AssignUser(todoID, userID)
}

func (s *TodoService) UnassignUser(todoID, userID uint) error {
	return s.repo.UnassignUser(todoID, userID)
}

func (s *TodoService) GetTodoAssignees(todoID uint) ([]models.User, error) {
	return s.repo.GetAssignees(todoID)
}

func (s *TodoService) GetTodoOwner(todoID uint) (*models.User, error) {
	return s.repo.GetOwner(todoID)
}

func (s *TodoService) ChangeOwner(todoID, userID uint) error {
	return s.repo.ChangeOwner(todoID, userID)
}
