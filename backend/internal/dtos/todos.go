package dtos

import (
	"time"
)

type TodoDTO struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	Tags        string     `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Owner       UserDTO    `json:"owner"`
	Assignees   []UserDTO  `json:"assignees"`
}
