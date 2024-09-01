package handlers

import (
	"time"
)

type TaskTemplateCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	// OwnerID     uint   `json:"owner_id" binding:"required"`
}

type TaskTemplateUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TaskTemplateResponse struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	OwnerID     uint         `json:"owner_id"`
	Owner       UserResponse `json:"owner"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type TaskTemplateListResponse struct {
	Templates []TaskTemplateResponse `json:"templates"`
	Total     int64                  `json:"total"`
}
