package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID          uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstName   string    `gorm:"type:text;not null" json:"first_name"`
	LastName    string    `gorm:"type:text;not null" json:"last_name"`
	Email       string    `gorm:"type:text;not null" json:"email"`
	PhoneNumber string    `gorm:"type:text;not null" json:"phone_number"`
	Address     string    `gorm:"type:text;not null" json:"address"`
	Date        time.Time `gorm:"type:time;not null" json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
}
