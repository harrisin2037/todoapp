package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/models"
)

type TaskTemplateRepository struct {
	db *gorm.DB
}

func NewTaskTemplateRepository(db *gorm.DB) *TaskTemplateRepository {
	return &TaskTemplateRepository{db: db}
}

func (r *TaskTemplateRepository) Create(template *models.TaskTemplate) error {
	return r.db.Create(template).Error
}

func (r *TaskTemplateRepository) GetList(sortBy, order string) ([]models.TaskTemplate, error) {
	var templates []models.TaskTemplate

	query := r.db.Preload("Owner")

	if sortBy != "" {
		if order == "desc" {
			query = query.Order(fmt.Sprintf("%s desc", sortBy))
		} else {
			query = query.Order(fmt.Sprintf("%s asc", sortBy))
		}
	}

	err := query.Find(&templates).Error

	return templates, err
}

func (r *TaskTemplateRepository) GetByID(id uint) (*models.TaskTemplate, error) {
	var template models.TaskTemplate
	err := r.db.Preload("Owner").First(&template, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

func (r *TaskTemplateRepository) Update(template *models.TaskTemplate) error {
	return r.db.Save(template).Error
}

func (r *TaskTemplateRepository) Delete(id uint) error {
	return r.db.Delete(&models.TaskTemplate{}, id).Error
}

func (r *TaskTemplateRepository) GetByOwnerID(ownerID uint) ([]models.TaskTemplate, error) {
	var templates []models.TaskTemplate
	err := r.db.Where("owner_id = ?", ownerID).Preload("Owner").Find(&templates).Error
	return templates, err
}

func (r *TaskTemplateRepository) GetByName(name string) (*models.TaskTemplate, error) {
	var template models.TaskTemplate
	err := r.db.Where("name = ?", name).First(&template).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}
