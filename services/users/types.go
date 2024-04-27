package users

import (
	"github.com/google/uuid"
	"github.com/razdacoder/mcwale-api/models"
)

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	UserExists(email string) (bool, error)
	GetUserByID(id string) (*models.User, error)
	CreateUser(payload RegisterUserPayload) error
	GetAllUsers() ([]models.User, error)
	UpdateUserInfo(user *models.User) error
	DeleteUser(id uuid.UUID) error
	UpdatePassword(id, password string) error
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=255"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type ResetPasswordPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordConfirmPayload struct {
	Password        string `json:"password" validate:"required,min=8,max=255"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=255"`
}
