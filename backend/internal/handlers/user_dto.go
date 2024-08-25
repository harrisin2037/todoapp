package handlers

import "github.com/harrisin2037/todoapp/internal/models"

type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserChangeRoleRequest struct {
	UserID  uint   `json:"user_id" binding:"required,min=1,numeric"`
	NewRole string `json:"new_role" binding:"required,oneof=user admin"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Color    string `json:"color"`
}

func NewUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role.String(),
		Color:    user.Color,
	}
}

func NewUsersResponse(users []models.User) []UserResponse {
	var usersResponse []UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, NewUserResponse(user))
	}
	return usersResponse
}
