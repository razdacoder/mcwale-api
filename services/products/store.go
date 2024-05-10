package products

import (
	"math"
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
	results := store.db.Select("Products").Delete(&category)
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

func (store *Store) GetAllProducts(category_slug string, minPrice, maxPrice float64, offset int, sortBy string) ([]models.Product, int, error) {
	var products []models.Product
	db := store.db.Model(&models.Product{})
	var total int64

	if minPrice != 0 {
		db = db.Where("price >= ?", minPrice)
	}
	if category_slug != "" {
		category, err := store.GetSingleCategory(category_slug)
		if err != nil {
			return nil, 0, err
		}
		db = db.Where("category_id = ?", category.ID)
	}
	if maxPrice != 0 {
		db = db.Where("price <= ?", maxPrice)
	}

	switch sortBy {
	case "new_arrivals":
		db = db.Order("created_at DESC")
	case "price_low_to_high":
		db = db.Order("price ASC")
	case "price_high_to_low":
		db = db.Order("price DESC")
	default:
		db = db.Order("created_at DESC")
	}
	db.Count(&total)

	pages := int(math.Ceil(float64(total) / float64(utils.ParseStringToInt(os.Getenv("PER_PAGE"), 10))))
	result := db.Offset(offset).Limit(utils.ParseStringToInt(os.Getenv("PER_PAGE"), 10)).Preload("Category").Find(&products)
	return products, pages, result.Error
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
	productMap := utils.ParseProductUpdate(product)
	results := store.db.Model(&existingProduct).Select("title", "images", "is_featured", "description", "style", "price", "discount_percentage").Updates(productMap)
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

func (store *Store) GetProductsByCategory(slug, style string, minPrice, maxPrice float64, offset int, sortBy string) ([]models.Product, error) {
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

	switch sortBy {
	case "new_arrivals":
		db = db.Order("created_at DESC")
	case "price_low_to_high":
		db = db.Order("price ASC")
	case "price_high_to_low":
		db = db.Order("price DESC")
	default:
		db = db.Order("created_at DESC")
	}

	results := db.Offset(offset).Limit(utils.ParseStringToInt(os.Getenv("PER_PAGE"), 10)).Preload("Category").Find(&products)
	return products, results.Error
}

func (store *Store) GetRecentProducts() ([]models.Product, error) {
	var products []models.Product
	results := store.db.Order("created_at desc").Limit(5).Preload("Category").Find(&products)
	return products, results.Error
}

func (store *Store) GetFeaturedProducts() ([]models.Product, error) {
	var products []models.Product
	results := store.db.Where("is_featured", true).Order("created_at desc").Limit(6).Preload("Category").Find(&products)
	return products, results.Error
}

func (store *Store) GetRelatedProducts(slug string) ([]models.Product, error) {
	var products []models.Product
	selectedProduct, err := store.GetSingleProduct(slug)
	if err != nil {
		return nil, err
	}
	results := store.db.Where("id != ?", selectedProduct.ID).Where("style = ?", selectedProduct.Style).Order("created_at desc").Limit(6).Preload("Category").Find(&products)

	return products, results.Error
}

func (store *Store) SearchProduct(query string, offset int, sortBy string) ([]models.Product, error) {
	var products []models.Product
	db := store.db.Model(&models.Product{})

	if query != "" {
		db = db.Where("(title ILIKE ? OR style ILIKE ?)",
			"%"+query+"%", "%"+query+"%")
	}

	switch sortBy {
	case "new_arrivals":
		db = db.Order("created_at DESC")
	case "price_low_to_high":
		db = db.Order("price ASC")
	case "price_high_to_low":
		db = db.Order("price DESC")
	default:
		db = db.Order("created_at DESC")
	}

	result := db.Offset(offset).Limit(utils.ParseStringToInt(os.Getenv("PER_PAGE"), 10)).Preload("Category").Find(&products)
	return products, result.Error
}
