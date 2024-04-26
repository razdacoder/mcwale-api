package users

import (
	"github.com/google/uuid"
	"github.com/razdacoder/mcwale-api/models"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	UserExists(email string) (bool, error)
	GetUserByID(id string) (*models.User, error)
	CreateUser(payload RegisterUserPayload) error
	GetAllUsers() ([]models.User, error)
	UpdateUserInfo(user *models.User) error
	DeleteUser(id uuid.UUID) error
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

func (store *Store) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := store.db.Model(&models.User{}).Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (store *Store) UserExists(email string) (bool, error) {
	var count int64
	result := store.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}

func (store *Store) GetUserByID(id string) (*models.User, error) {
	var user models.User
	result := store.db.Model(&models.User{}).Where("id = ?", id).First(&user)
	return &user, result.Error
}

func (store *Store) CreateUser(payload RegisterUserPayload) error {
	user := &models.User{
		Firstname: payload.FirstName,
		Lastname:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	}
	result := store.db.Create(user)
	return result.Error
}

func (store *Store) GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := store.db.Find(&users)
	return users, result.Error
}

func (store *Store) UpdateUserInfo(user *models.User) error {
	results := store.db.Save(&user)
	return results.Error
}

func (store *Store) DeleteUser(id uuid.UUID) error {
	results := store.db.Delete(&models.User{}, id)
	return results.Error
}
