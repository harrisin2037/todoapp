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

func (s *TodoService) CreateTodo(todo *models.Todo) error {
	return s.repo.Create(todo)
}

func (s *TodoService) GetTodos(statuses []string, sortBy, order string) ([]models.Todo, error) {
	return s.repo.GetList(statuses, sortBy, order)
}

func (s *TodoService) GetTodo(id int64) (*models.Todo, error) {
	return s.repo.GetByID(uint64(id))
}

func (s *TodoService) UpdateTodo(todo *models.Todo) error {
	return s.repo.Update(todo)
}

func (s *TodoService) DeleteTodo(id int64) error {
	return s.repo.Delete(uint64(id))
}
