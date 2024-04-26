package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	Admin    UserRole = "admin"
	Customer UserRole = "customer"
)

func (e *UserRole) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = UserRole(str)

	return nil
}

func (e UserRole) Value() (interface{}, error) {
	return string(e), nil
}

type User struct {
	ID        uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Firstname string     `gorm:"type:text;not null" json:"first_name"`
	Lastname  string     `gorm:"type:text;not null" json:"last_name"`
	Email     string     `gorm:"type:text;unique;not null" json:"email"`
	UserRole  *UserRole  `gorm:"type:user_role;default:'customer';not null" json:"role"`
	Password  string     `gorm:"type:text;not null" json:"-"`
	CreatedAt time.Time  `gorm:"auto_now_add" json:"created_at"`
	UpdatedAt time.Time  `gorm:"auto_now" json:"updated_at"`
	DeletedAt *time.Time `gorm:"auto_now_add;default:null" json:"-"`
}
