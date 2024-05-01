package orders

import (
	"github.com/google/uuid"
	"github.com/razdacoder/mcwale-api/models"
)

type OrderStore interface {
	CreateOrder(CreateOrderPayload) error
	GetOrders() ([]models.Order, error)
}

type CreateOrderPayload struct {
	OrderNumber string                   `json:"order_number" validate:"required"`
	FirstName   string                   `json:"first_name" validate:"required"`
	LastName    string                   `json:"last_name" validate:"required"`
	Email       string                   `json:"email" validate:"required"`
	PhoneNumber string                   `json:"phone_number" validate:"required"`
	Address1    string                   `json:"address_line_1" validate:"required"`
	Address2    string                   `json:"address_line_2"`
	Town        string                   `json:"town" validate:"required"`
	State       string                   `json:"state" validate:"required"`
	Country     string                   `json:"country" validate:"required"`
	PostalCode  string                   `json:"postal_code" validate:"required"`
	OrderNote   string                   `json:"order_note"`
	Total       float64                  `json:"total" validate:"required"`
	Items       []CreateOrderItemPayload `json:"items" validate:"required"`
}

type CreateOrderItemPayload struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Size      string    `json:"size" validate:"required"`
	Color     string    `json:"color" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required;min=1"`
}

type UpdateOrderPayload struct {
	Status string `json:"status" validate:"required"`
}
