package orders

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

func (store *Store) CreateOrder(payload CreateOrderPayload) (string, error) {
	order := &models.Order{
		ID:          uuid.New(),
		OrderNumber: payload.OrderNumber,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Address1:    payload.Address1,
		Address2:    payload.Address2,
		Town:        payload.Town,
		State:       payload.State,
		Country:     payload.Country,
		PostalCode:  payload.PostalCode,
		OrderNote:   payload.OrderNote,
		Total:       payload.Total,
	}

	for _, itemPayload := range payload.Items {
		item := &models.OrderItem{
			ID:        uuid.New(),
			ProductID: itemPayload.ProductID,
			OrderID:   order.ID,
			Size:      itemPayload.Size,
			Color:     itemPayload.Color,
			Quantity:  itemPayload.Quantity,
		}
		order.Items = append(order.Items, *item)
	}

	tx := store.db.Begin()
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Model(order).Association("Items").Append(order.Items); err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return order.OrderNumber, nil
}

func (store *Store) GetOrders() ([]models.Order, error) {
	var orders []models.Order
	results := store.db.Find(&orders)
	if results.Error != nil {
		return nil, results.Error
	}

	return orders, nil
}

func (store *Store) GetOrderByID(id string) (*models.Order, error) {
	var order models.Order
	results := store.db.Model(&models.Order{}).Where("id = ?", id).Preload("Items.Product.Category").First(&order)
	if results.Error != nil {
		return nil, results.Error
	}

	return &order, nil
}

func (store *Store) UpdateOrderStatus(id, status string) error {
	results := store.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status)
	return results.Error
}

func (store *Store) DeleteOrder(id string) error {
	order, err := store.GetOrderByID(id)
	if err != nil {
		return err
	}
	results := store.db.Select("Items").Delete(&order)
	return results.Error
}
