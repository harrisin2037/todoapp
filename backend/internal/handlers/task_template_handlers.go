package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/models"
	"github.com/harrisin2037/todoapp/internal/service"
	"github.com/harrisin2037/todoapp/internal/websocket"
)

type TaskTemplateHandler struct {
	hub         *websocket.Hub
	userService *service.UserService
	service     *service.TaskTemplateService
}

func NewTaskTemplateHandler(service *service.TaskTemplateService, userService *service.UserService, hub *websocket.Hub) *TaskTemplateHandler {
	return &TaskTemplateHandler{
		service:     service,
		userService: userService,
		hub:         hub,
	}
}

func (h *TaskTemplateHandler) CreateTaskTemplate(c *gin.Context) {

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

	var req TaskTemplateCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	owner, err := h.userService.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner ID"})
		return
	}

	taskTemplate := &models.TaskTemplate{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     owner.ID,
		Owner:       *owner,
	}

	if err := h.service.CreateTaskTemplate(taskTemplate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast <- []byte(`{"message": "new task template created"}`)

	response := TaskTemplateResponse{
		ID:          taskTemplate.ID,
		Name:        taskTemplate.Name,
		Description: taskTemplate.Description,
		OwnerID:     taskTemplate.OwnerID,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TaskTemplateHandler) GetTaskTemplates(c *gin.Context) {

	claims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not recognized"})
		return
	}

	user, ok := claims.(*models.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	var (
		err       error
		response  = []TaskTemplateResponse{}
		templates = []models.TaskTemplate{}
	)

	if user.Role == "admin" {
		templates, err = h.service.GetTaskTemplateList("id", "asc")
	} else {
		templates, err = h.service.GetTaskTemplatesByOwnerID(user.UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, template := range templates {
		response = append(response, TaskTemplateResponse{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			OwnerID:     template.OwnerID,
			Owner:       NewUserResponse(template.Owner),
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskTemplateHandler) GetTaskTemplatesByOwnerID(c *gin.Context) {

	ownerID, err := strconv.ParseUint(c.Param("ownerID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid owner ID"})
		return
	}

	templates, err := h.service.GetTaskTemplatesByOwnerID(uint(ownerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task templates not found"})
		return
	}

	response := []TaskTemplateResponse{}
	for _, template := range templates {
		response = append(response, TaskTemplateResponse{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			OwnerID:     template.OwnerID,
			Owner:       NewUserResponse(template.Owner),
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskTemplateHandler) GetTaskTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	template, err := h.service.GetTaskTemplateByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task template not found"})
		return
	}

	response := TaskTemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		OwnerID:     template.OwnerID,
		Owner:       NewUserResponse(template.Owner),
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskTemplateHandler) UpdateTaskTemplate(c *gin.Context) {

	userClaims, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not recognized"})
		return
	}

	claims, ok := userClaims.(*models.Claims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req TaskTemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := h.service.GetTaskTemplateByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task template not found"})
		return
	}

	if template.OwnerID != claims.UserID && claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this task template"})
		return
	}

	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}

	if err := h.service.UpdateTaskTemplate(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast <- []byte(fmt.Sprintf(`
		{
			"message":"task template updated",
			"taskTemplate":{
				"id":%d,
				"name":"%s",
				"description":"%s",
				"owner_id":%d
			}
		}`, template.ID, template.Name, template.Description, template.OwnerID))

	response := TaskTemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		OwnerID:     template.OwnerID,
		Owner:       NewUserResponse(template.Owner),
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskTemplateHandler) DeleteTaskTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteTaskTemplate(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task template not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
