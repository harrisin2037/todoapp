package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/models"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *TodoRepository) GetList(statuses []string, sortBy, order string) ([]models.Todo, error) {

	var todos []models.Todo

	query := r.db

	if len(statuses) > 0 {
		query = query.Where("status IN (?)", statuses)
	}

	if sortBy != "" {
		if order == "desc" {
			query = query.Order(fmt.Sprintf("%s desc", sortBy))
		} else {
			query = query.Order(fmt.Sprintf("%s asc", sortBy))
		}
	}

	err := query.Find(&todos).Error

	return todos, err
}

func (r *TodoRepository) GetByID(id uint64) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *TodoRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Todo{}, id).Error
}
