package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Category struct {
	ID        uuid.UUID      `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title     string         `gorm:"type:text;not null" json:"title"`
	Slug      string         `gorm:"type:text;unique;not null" json:"slug"`
	Styles    pq.StringArray `gorm:"type:text[];not null" json:"styles"`
	Image     string         `gorm:"type:string;not null" json:"image"`
	Products  []Product      `json:"products,omitempty"`
	CreatedAt time.Time      `gorm:"auto_now_add" json:"-"`
	UpdatedAt time.Time      `gorm:"auto_now" json:"-"`
	DeletedAt *time.Time     `gorm:"auto_now_add;default:null" json:"-"`
}

type Product struct {
	ID                 uuid.UUID      `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title              string         `gorm:"type:text;not null" json:"title"`
	Slug               string         `gorm:"type:text;unique;not null" json:"slug"`
	Images             pq.StringArray `gorm:"type:text[];not null" json:"images"`
	Style              string         `gorm:"type:string;not null" json:"style"`
	IsFeatured         bool           `gorm:"default:false" json:"is_featured"`
	Price              float64        `gorm:"type:decimal(10, 2)" json:"price"`
	Description        string         `gorm:"type:text;not null" json:"description"`
	DiscountPercentage float64        `gorm:"type:decimal(5, 2)" json:"discount_percentage"`
	CategoryID         uuid.UUID      `gorm:"index" json:"-"`
	Category           Category       `json:"category,omitempty"`
	CreatedAt          time.Time      `gorm:"auto_now_add" json:"-"`
	UpdatedAt          time.Time      `gorm:"auto_now" json:"-"`
	DeletedAt          *time.Time     `gorm:"auto_now_add;default:null" json:"-"`
}
