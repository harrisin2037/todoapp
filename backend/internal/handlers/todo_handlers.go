package handlers

import (
	"fmt"
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
	hub         *websocket.Hub
	userService *service.UserService
	service     *service.TodoService
}

func NewTodoHandler(service *service.TodoService, userService *service.UserService, hub *websocket.Hub) *TodoHandler {
	return &TodoHandler{
		service:     service,
		userService: userService,
		hub:         hub,
	}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {

	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	claims, ok := userClaims.(*models.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

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

	owner, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner ID"})
		return
	}

	assignees := []models.User{}
	if len(req.AssigneeIDs) > 0 {
		assignees, err = h.userService.GetUsersByIDs(req.AssigneeIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
			return
		}
	}

	todo := &models.Todo{
		Name:        req.Name,
		Description: req.Description,
		DueDate:     dueDate,
		Status:      req.Status,
		OwnerID:     owner.ID,
		Owner:       *owner,
		Assignees:   assignees,
	}

	if todo.Status == "" {
		todo.Status = "pending"
	}

	if err := h.service.CreateTodo(todo, req.AssigneeIDs); err != nil {
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
		OwnerID:     todo.OwnerID,
		Owner:       NewUserResponse(todo.Owner),
		Assignees:   NewUsersResponse(todo.Assignees),
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TodoHandler) GetTodos(c *gin.Context) {

	var (
		statusQuery = c.Query("status")
		sortBy      = c.Query("sort_by")
		order       = c.Query("order")
	)

	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not regconized"})
		return
	}

	user, ok := claims.(*models.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

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

	var (
		err      error
		response = []TodoResponse{}
		todos    = []models.Todo{}
	)

	if user.Role == "admin" {
		todos, err = h.service.GetTodosByAdmin(statuses, sortBy, order)
	} else {
		todos, err = h.service.GetTodos(statuses, sortBy, order, user.UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, todo := range todos {
		response = append(response, TodoResponse{
			ID:          todo.ID,
			Name:        todo.Name,
			Description: todo.Description,
			DueDate:     todo.DueDate,
			Status:      todo.Status,
			OwnerID:     todo.OwnerID,
			Owner:       NewUserResponse(todo.Owner),
			Assignees:   NewUsersResponse(todo.Assignees),
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

	todo, err := h.service.GetTodo(uint(id))
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
		OwnerID:     todo.OwnerID,
		Owner:       NewUserResponse(todo.Owner),
		Assignees:   NewUsersResponse(todo.Assignees),
	}

	c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {

	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not regconized"})
		return
	}

	claims, ok := userClaims.(*models.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userId := claims.UserID
	user, err := h.userService.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

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

	todo, err := h.service.GetTodo(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if todo.OwnerID != userId && user.Role != "admin" && !models.AssigneesContainsUser(todo.Assignees, *user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this todo"})
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
		todo.DueDate = &parsedTime
	}
	if req.Status != "" {
		todo.Status = req.Status
	}
	if req.OwnerID != nil && *req.OwnerID == 0 && todo.OwnerID == 0 && user.Role == "admin" {
		err := h.service.ChangeOwner(todo.ID, user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "admin change owner error"})
			return
		}
	}
	if req.OwnerID != nil && *req.OwnerID > 0 {
		if user.Role == "admin" {
			newOwner, err := h.userService.GetUserByID(*req.OwnerID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid new owner ID"})
				return
			}
			todo.Owner = *newOwner
			todo.OwnerID = *req.OwnerID
		}
	}
	if len(req.AssigneeIDs) > 0 {
		assignees, err := h.userService.GetUsersByIDs(req.AssigneeIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignee ID"})
			return
		}
		todo.Assignees = assignees
	}

	if err := h.service.UpdateTodo(todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.hub.Broadcast <- []byte(
		fmt.Sprintf(`
		{
			"message":"todo updated",
			"todo":{
				"id":%d,
				"name":"%s",
				"description":"%s",
				"due_date":"%s",
				"status":"%s",
				"owner_id":%d
			}
		}`,
			todo.ID, todo.Name, todo.Description, todo.DueDate.Format(time.RFC3339), todo.Status, todo.OwnerID))

	response := TodoResponse{
		ID:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
		DueDate:     todo.DueDate,
		Status:      todo.Status,
		OwnerID:     todo.OwnerID,
		Owner:       NewUserResponse(todo.Owner),
		Assignees:   NewUsersResponse(todo.Assignees),
	}

	c.JSON(http.StatusOK, response)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteTodo(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
