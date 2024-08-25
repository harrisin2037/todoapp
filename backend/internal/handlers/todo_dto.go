package handlers

import (
	"time"
)

type TodoCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date"`
	Status      string  `json:"status"`
	AssigneeIDs []uint  `json:"assignee_ids"`
}

type TodoUpdateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	DueDate     *string `json:"due_date"`
	Status      string  `json:"status"`
	OwnerID     *uint   `json:"owner_id"`
	AssigneeIDs []uint  `json:"assignee_ids"`
}

type TodoResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	DueDate     *time.Time     `json:"due_date"`
	Status      string         `json:"status"`
	OwnerID     uint           `json:"owner_id"`
	Owner       UserResponse   `json:"owner"`
	Assignees   []UserResponse `json:"assignees"`
}
