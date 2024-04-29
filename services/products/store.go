package products

import (
	"os"

	"github.com/razdacoder/mcwale-api/models"
	"github.com/razdacoder/mcwale-api/utils"
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

func (store *Store) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	result := store.db.Find(&categories)
	return categories, result.Error
}

func (store *Store) CreateCategory(payload CreateCategoryPayload) error {
	category := &models.Category{
		Title:  payload.Title,
		Slug:   payload.Slug,
		Styles: payload.Styles,
		Image:  payload.Image,
	}
	result := store.db.Create(category)
	return result.Error
}

func (store *Store) GetSingleCategory(slug string) (*models.Category, error) {
	var category models.Category
	result := store.db.Model(&models.Category{}).Where("slug = ?", slug).First(&category)
	return &category, result.Error
}

func (store *Store) UpdateCategory(slug string, category *models.Category) error {
	existingCategory, err := store.GetSingleCategory(slug)
	if err != nil {
		return err
	}
	results := store.db.Model(&existingCategory).Updates(&category)
	return results.Error
}

func (store *Store) DeleteCategory(slug string) error {
	category, err := store.GetSingleCategory(slug)
	if err != nil {
		return err
	}
	results := store.db.Delete(&models.Category{}, category.ID)
	return results.Error
}

func (store *Store) CreateProduct(payload CreateProductPayload) error {
	product := &models.Product{
		Title:              payload.Title,
		Slug:               payload.Slug,
		Images:             payload.Images,
		Style:              payload.Style,
		IsFeatured:         payload.IsFeatured,
		Price:              payload.Price,
		Description:        payload.Description,
		DiscountPercentage: payload.DiscountPercentage,
		CategoryID:         payload.CategoryID,
	}

	result := store.db.Create(product)
	return result.Error
}

func (store *Store) GetAllProducts(style, category_slug string, minPrice, maxPrice float64, offset int) ([]models.Product, error) {
	var products []models.Product
	db := store.db.Model(&models.Product{})
	if style != "" {
		db = db.Where("style = ?", style)
	}
	if minPrice != 0 {
		db = db.Where("price >= ?", minPrice)
	}
	if category_slug != "" {
		category, err := store.GetSingleCategory(category_slug)
		if err != nil {
			return nil, err
		}
		db = db.Where("category_id = ?", category.ID)
	}
	if maxPrice != 0 {
		db = db.Where("price <= ?", maxPrice)
	}
	result := db.Offset(offset).Limit(1).Preload("Category").Find(&products)
	return products, result.Error
}

func (store *Store) GetSingleProduct(slug string) (*models.Product, error) {
	var product models.Product
	result := store.db.Model(&models.Product{}).Where("slug = ?", slug).Preload("Category").First(&product)
	return &product, result.Error
}

func (store *Store) UpdateProduct(slug string, product *models.Product) error {
	existingProduct, err := store.GetSingleProduct(slug)
	if err != nil {
		return err
	}
	results := store.db.Model(&existingProduct).Updates(&product)
	return results.Error
}

func (store *Store) DeleteProduct(slug string) error {
	product, err := store.GetSingleProduct(slug)
	if err != nil {
		return err
	}
	results := store.db.Delete(&models.Product{}, product.ID)
	return results.Error
}

func (store *Store) GetProductsByCategory(slug, style string, minPrice, maxPrice float64, offset int) ([]models.Product, error) {
	var products []models.Product
	category, err := store.GetSingleCategory(slug)
	if err != nil {
		return nil, err
	}
	db := store.db.Model(&models.Product{}).Where("category_id = ?", category.ID)
	if style != "" {
		db = db.Where("style = ?", style)
	}
	if minPrice != 0 {
		db = db.Where("price >= ?", minPrice)
	}

	if maxPrice != 0 {
		db = db.Where("price <= ?", maxPrice)
	}

	results := db.Offset(offset).Limit(utils.ParseStringToInt(os.Getenv("PER_PAGE"), 10)).Preload("Category").Find(&products)
	return products, results.Error
}
