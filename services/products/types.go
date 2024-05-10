package products

import (
	"time"

	"github.com/google/uuid"
	"github.com/razdacoder/mcwale-api/models"
)

type ProductStore interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(CreateCategoryPayload) error
	GetSingleCategory(string) (*models.Category, error)
	UpdateCategory(*models.Category) error
	DeleteCategory(string) error
	CreateProduct(CreateProductPayload) error
	GetAllProducts(string, string, float64, float64, int) ([]models.Product, error)
	GetSingleProduct(string) (*models.Product, error)
	UpdateProduct(string, *models.Product) error
	DeleteProduct(string) error
	GetProductsByCategory(string, string, float64, float64, int) ([]models.Product, error)
}

type CreateCategoryPayload struct {
	Title  string   `json:"title" validate:"required"`
	Slug   string   `json:"slug" validate:"required"`
	Styles []string `json:"styles" validate:"required"`
	Image  string   `json:"image" validate:"required"`
}

type CreateProductPayload struct {
	Title              string    `json:"title" validate:"required"`
	Slug               string    `json:"slug" validate:"required"`
	Images             []string  `json:"images" validate:"required"`
	Style              string    `json:"style" validate:"required"`
	IsFeatured         bool      `json:"is_featured" validate:"required"`
	Price              float64   `json:"price" validate:"required"`
	Description        string    `json:"description" validate:"required"`
	DiscountPercentage float64   `json:"discount_percentage"`
	CategoryID         uuid.UUID `json:"category_id" validate:"required"`
}

type CategoryReturn struct {
	ID        uuid.UUID        `json:"id"`
	Title     string           `json:"title"`
	Slug      string           `json:"slug"`
	Styles    []string         `json:"styles"`
	Image     string           `json:"image"`
	Products  []models.Product `json:"-"`
	CreatedAt time.Time        `json:"-"`
	UpdatedAt time.Time        `json:"-"`
	DeletedAt *time.Time       `json:"-"`
}
