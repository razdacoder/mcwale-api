package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Shipped   OrderStatus = "shipped"
	Delivered OrderStatus = "delivered"
)

func (e *OrderStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid str")
	}
	*e = OrderStatus(str)

	return nil
}

func (e OrderStatus) Value() (interface{}, error) {
	return string(e), nil
}

type Order struct {
	ID          uuid.UUID    `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderNumber string       `gorm:"type:text;unique;not null" json:"order_number"`
	FirstName   string       `gorm:"type:text;not null" json:"first_name"`
	LastName    string       `gorm:"type:text;not null" json:"last_name"`
	Email       string       `gorm:"type:text;not null" json:"email"`
	PhoneNumber string       `gorm:"type:text;not null" json:"phone_number"`
	Address1    string       `gorm:"type:text;not null" json:"address_line_1"`
	Address2    string       `gorm:"type:text" json:"address_line_2"`
	Town        string       `gorm:"type:text;not null" json:"town"`
	State       string       `gorm:"type:text;not null" json:"state"`
	Country     string       `gorm:"type:text;not null" json:"country"`
	PostalCode  string       `gorm:"type:text;not null" json:"postal_code"`
	OrderNote   string       `gorm:"type:text;not null" json:"order_note"`
	Status      *OrderStatus `gorm:"type:order_status;default:'pending';not null" json:"status"`
	Total       float64      `gorm:"type:decimal(10, 2);not null" json:"total"`
	Items       []OrderItem  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"-"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	ProductID uuid.UUID `gorm:"index;not null" json:"-"`
	Product   Product   `json:"product,omitempty"`
	OrderID   uuid.UUID `gorm:"index;not null" json:"-"`
	Order     Order     `json:"-"`
	Size      string    `gorm:"type:text;not null" json:"size"`
	Color     string    `gorm:"type:text;not null" json:"color"`
	Quantity  int       `json:"quantity"`
}
