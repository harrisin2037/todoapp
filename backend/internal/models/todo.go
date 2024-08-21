package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date" gorm:"default:null"`
	Status      string     `json:"status" gorm:"default:'pending'"`
	Priority    string     `json:"priority"`
	Tags        string     `json:"tags"`
}
