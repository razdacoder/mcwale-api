package appointments

import (
	"fmt"

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

func (store *Store) GetAllAppointments() ([]models.Appointment, error) {
	var appointments []models.Appointment
	result := store.db.Find(&appointments)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return appointments, nil
}

func (store *Store) CreateNewAppointments(app CreateAppointmentPayload) error {
	appointment := &models.Appointment{
		FirstName:   app.FirstName,
		LastName:    app.LastName,
		Email:       app.Email,
		PhoneNumber: app.Email,
		Address:     app.Address,
		Date:        app.Date,
	}

	if err := store.db.Create(&appointment).Error; err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (store *Store) GetSingleAppointment(id string) (*models.Appointment, error) {
	var appointment *models.Appointment

	if err := store.db.Where("id = ?", id).First(&appointment).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return appointment, nil
}
