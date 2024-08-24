package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/service"
	"github.com/harrisin2037/todoapp/internal/websocket"
)

type TodoHandler struct {
	hub     *websocket.Hub
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService, hub *websocket.Hub) *TodoHandler {
	return &TodoHandler{service: service, hub: hub}
}

type TodoCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date"`
	Status      string  `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
}

type TodoUpdateRequest struct {
	Name        string  `json:"name" binding:"omitempty"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date"`
	Status      string  `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
}

type TodoResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status"`
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {

	var req TodoCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		parsedTime, err := dateparse.ParseAny(*req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		dueDate = &parsedTime
	}

	// if dueDate != nil && dueDate.Before(time.Now()) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Due date must be in the future"})
	// 	return
	// }

	todo := &models.Todo{
		Name:        req.Name,
		Description: req.Description,
		DueDate:     dueDate,
		Status:      req.Status,
	}

	if todo.Status == "" {
		todo.Status = "pending"
	}

	if err := h.service.CreateTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast <- []byte(`{"message": "new todo created"}`)

	response := TodoResponse{
		ID:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
		DueDate:     todo.DueDate,
		Status:      todo.Status,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	var (
		statusQuery = c.Query("status")
		sortBy      = c.Query("sort_by")
		order       = c.Query("order")
	)

	statuses := []string{}
	if statusQuery != "" {
		statuses = strings.Split(strings.ToLower(statusQuery), ",")
		for _, status := range statuses {
			if status != "pending" && status != "in_progress" && status != "completed" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status field"})
				return
			}
		}
	}

	if sortBy != "" && sortBy != "name" && sortBy != "due_date" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort_by field"})
		return
	}
	if order != "" && order != "asc" && order != "desc" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order field"})
		return
	}

	todos, err := h.service.GetTodos(statuses, sortBy, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := []TodoResponse{}

	for _, todo := range todos {
		response = append(response, TodoResponse{
			ID:          todo.ID,
			Name:        todo.Name,
			Description: todo.Description,
			DueDate:     todo.DueDate,
			Status:      todo.Status,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	todo, err := h.service.GetTodo(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	response := TodoResponse{
		ID:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
		DueDate:     todo.DueDate,
		Status:      todo.Status,
	}

	c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req TodoUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.service.GetTodo(int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if req.Name != "" {
		todo.Name = req.Name
	}
	if req.Description != "" {
		todo.Description = req.Description
	}
	if req.DueDate != nil {
		parsedTime, err := dateparse.ParseAny(*req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		// if parsedTime.Before(time.Now()) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Due date must be in the future"})
		// 	return
		// }
		todo.DueDate = &parsedTime
	}
	if req.Status != "" {
		todo.Status = req.Status
	}

	if err := h.service.UpdateTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast <- []byte(`{"message": "todo updated"}`)

	response := TodoResponse{
		ID:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
		DueDate:     todo.DueDate,
		Status:      todo.Status,
	}

	c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteTodo(int64(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
