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

func (r *TodoRepository) Create(todo *models.Todo, assigneeIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(todo).Error; err != nil {
			return err
		}

		if len(assigneeIDs) > 0 {
			var assignees []models.User
			if err := tx.Where("id IN ?", assigneeIDs).Find(&assignees).Error; err != nil {
				return err
			}
			return tx.Model(todo).Association("Assignees").Append(assignees)
		}

		return nil
	})
}

func (r *TodoRepository) GetListByAdmin(statuses []string, sortBy, order string) ([]models.Todo, error) {
	var todos []models.Todo

	query := r.db.Preload("Owner").Preload("Assignees")

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

func (r *TodoRepository) GetList(statuses []string, sortBy, order string, userID uint) ([]models.Todo, error) {

	var todos []models.Todo

	query := r.db.Preload("Owner").Preload("Assignees")

	if len(statuses) > 0 {
		query = query.Where("status IN (?)", statuses)
	}

	if userID != 0 {
		query = query.Where("owner_id = ?", userID)
		query = query.Or("id IN (SELECT todo_id FROM todo_assignees WHERE user_id = ?)", userID)
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

func (r *TodoRepository) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Preload("Owner").Preload("Assignees").First(&todo, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepository) Update(todo *models.Todo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(todo).Error; err != nil {
			return err
		}

		if err := tx.Model(todo).Association("Owner").Replace(&todo.Owner); err != nil {
			return err
		}

		unassignedTodo := &models.Todo{Model: gorm.Model{ID: todo.ID}}

		if err := tx.Model(&unassignedTodo).Association("Assignees").Clear(); err != nil {
			return err
		}

		if len(todo.Assignees) > 0 {
			return tx.Model(todo).Association("Assignees").Append(todo.Assignees)
		}

		return nil
	})
}

func (r *TodoRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		todo := &models.Todo{Model: gorm.Model{ID: id}}

		if err := tx.Model(todo).Association("Assignees").Clear(); err != nil {
			return err
		}

		return tx.Delete(todo).Error
	})
}

func (r *TodoRepository) AssignUser(todoID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		todo := &models.Todo{Model: gorm.Model{ID: todoID}}
		user := &models.User{Model: gorm.Model{ID: userID}}

		if err := tx.First(todo).Error; err != nil {
			return err
		}

		if err := tx.First(user).Error; err != nil {
			return err
		}

		return tx.Model(todo).Association("Assignees").Append(user)
	})
}

func (r *TodoRepository) UnassignUser(todoID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		todo := &models.Todo{Model: gorm.Model{ID: todoID}}
		user := &models.User{Model: gorm.Model{ID: userID}}

		if err := tx.First(todo).Error; err != nil {
			return err
		}

		if err := tx.First(user).Error; err != nil {
			return err
		}

		return tx.Model(todo).Association("Assignees").Delete(user)
	})
}

func (r *TodoRepository) GetAssignees(todoID uint) ([]models.User, error) {
	var users []models.User
	err := r.db.Model(&models.Todo{Model: gorm.Model{ID: todoID}}).Association("Assignees").Find(&users)
	return users, err
}

func (r *TodoRepository) GetOwner(todoID uint) (*models.User, error) {
	var user models.User
	err := r.db.Model(&models.Todo{Model: gorm.Model{ID: todoID}}).Association("Owner").Find(&user)
	return &user, err
}

func (r *TodoRepository) ChangeOwner(todoID, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		todo := &models.Todo{Model: gorm.Model{ID: todoID}}
		user := &models.User{Model: gorm.Model{ID: userID}}

		if err := tx.First(todo).Error; err != nil {
			return err
		}

		if err := tx.First(user).Error; err != nil {
			return err
		}

		return tx.Model(todo).Association("Owner").Replace(user)
	})
}
