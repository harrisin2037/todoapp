package models

import (
	"gorm.io/gorm"
)

type TaskTemplate struct {
	gorm.Model
	Name                     string `json:"name" gorm:"type:varchar(255);not null"`
	Description              string `json:"description" gorm:"type:text"`
	DefaultDurationTimestamp int    `json:"default_duration" gorm:"not null"`
	OwnerID                  uint   `json:"owner_id" gorm:"not null"`
	Owner                    User   `json:"owner" gorm:"foreignKey:OwnerID"`
}
