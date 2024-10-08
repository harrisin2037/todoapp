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

func (r *UserRepository) GetUsersByIDs(ids []uint) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("id IN ?", ids).Find(&users).Error
	return users, err
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

func (r *UserRepository) UpdateRole(userID uint, role models.Role) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("role", role).Error
}

func (r *UserRepository) FindByID(userID uint) (*models.User, error) {
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

func (r *UserRepository) Update(user *models.User) error {

	var existingUser models.User

	err := r.db.First(&existingUser, user.ID).Error
	if err != nil {
		return err
	}

	if existingUser.Password != user.Password {
		err = r.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(user).Error
	} else {
		err = r.db.Model(&models.User{}).Where("id = ?", user.ID).Omit("password").Updates(user).Error
	}

	return err
}

func (r *UserRepository) Delete(userID uint) error {
	return r.db.Delete(&models.User{}, userID).Error
}
