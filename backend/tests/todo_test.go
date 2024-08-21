package tests

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/repository"
	"github.com/harrisin2037/todoapp/internal/service"
)

func TestCreateTodo(t *testing.T) {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.Todo{})

	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)

	todo := &models.Todo{
		Name:        "Test Todo",
		Description: "This is a test todo",
	}

	err = todoService.CreateTodo(todo)
	if err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	var createdTodo models.Todo
	result := db.First(&createdTodo, "name = ?", "Test Todo")
	if result.Error != nil {
		t.Fatalf("Failed to find created todo: %v", result.Error)
	}

	if createdTodo.Description != "This is a test todo" {
		t.Errorf("Expected description to be 'This is a test todo', got '%s'", createdTodo.Description)
	}
}
