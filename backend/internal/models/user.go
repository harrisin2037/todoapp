package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/utils"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var (
	ErrInvalidRole = gorm.ErrInvalidData
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"type:varchar(20);default:'user'"`
	Color    string `gorm:"type:varchar(30);default:'#000000'"`
}

func (r Role) String() string {
	return string(r)
}

func ParseRole(role string) (Role, error) {
	switch role {
	case "user":
		return RoleUser, nil
	case "admin":
		return RoleAdmin, nil
	default:
		return "", ErrInvalidRole
	}
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	u.Color = utils.GenerateColor(u.ID)
	return tx.Save(u).Error
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	if u.Color == "" {
		u.Color = utils.GenerateColor(u.ID)
		return tx.Save(u).Error
	}
	return nil
}

func AssigneesContainsUser(assignees []User, user User) bool {
	for _, assignee := range assignees {
		if assignee.ID == user.ID {
			return true
		}
	}
	return false
}
