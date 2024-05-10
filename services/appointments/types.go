package appointments

import "time"

type CreateAppointmentPayload struct {
	FirstName   string    `json:"first_name" validate:"required"`
	LastName    string    `json:"last_name" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	PhoneNumber string    `json:"phone_number" validate:"required"`
	Address     string    `json:"address" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
}
