package service

import (
	"errors"
	"time"

	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/repository"
)

var ErrTemplateNotFound = errors.New("template not found")

type TaskTemplateService struct {
	repo     *repository.TaskTemplateRepository
	todoRepo *repository.TodoRepository
}

func NewTaskTemplateService(repo *repository.TaskTemplateRepository) *TaskTemplateService {
	return &TaskTemplateService{repo: repo}
}

func (s *TaskTemplateService) CreateTaskTemplate(template *models.TaskTemplate) error {
	return s.repo.Create(template)
}

func (s *TaskTemplateService) GetTaskTemplateList(sortBy, order string) ([]models.TaskTemplate, error) {
	return s.repo.GetList(sortBy, order)
}

func (s *TaskTemplateService) GetTaskTemplateByID(id uint) (*models.TaskTemplate, error) {
	return s.repo.GetByID(id)
}

func (s *TaskTemplateService) UpdateTaskTemplate(template *models.TaskTemplate) error {
	return s.repo.Update(template)
}

func (s *TaskTemplateService) DeleteTaskTemplate(id uint) error {
	return s.repo.Delete(id)
}

func (s *TaskTemplateService) GetTaskTemplatesByOwnerID(ownerID uint) ([]models.TaskTemplate, error) {
	return s.repo.GetByOwnerID(ownerID)
}

func (s *TaskTemplateService) GetTaskTemplateByName(name string) (*models.TaskTemplate, error) {
	return s.repo.GetByName(name)
}

func (s *TaskTemplateService) CreateTaskFromTemplate(templateID uint, ownerID uint) (*models.Todo, error) {

	template, err := s.repo.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	if template == nil {
		return nil, ErrTemplateNotFound
	}

	if template.OwnerID == 0 {
		template.OwnerID = ownerID
	}

	dueDate := time.Now()
	if template.DefaultDurationTimestamp != 0 {
		dueDate = time.Now().Add(time.Duration(template.DefaultDurationTimestamp) * time.Second)
	}

	todo := &models.Todo{
		Name:        template.Name,
		Description: template.Description,
		Status:      "pending",
		DueDate:     &dueDate,
		OwnerID:     ownerID,
		Owner:       template.Owner,
	}

	err = s.todoRepo.Create(todo, []uint{ownerID})
	if err != nil {
		return nil, err
	}

	return todo, nil
}
