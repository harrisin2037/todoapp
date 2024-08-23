package repository

import (
	"gorm.io/gorm"

	"github.com/harrisin2037/todoapp/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateRole(userID int64, role models.Role) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func (r *UserRepository) FindByID(userID int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(userID int64, username, email, password string) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err
	}

	if username != "" {
		user.Username = username
	}

	if email != "" {
		user.Email = email
	}

	if password != "" {
		err = user.SetPassword(password)
		if err != nil {
			return err
		}
	}

	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(userID int64) error {
	return r.db.Delete(&models.User{}, userID).Error
}
