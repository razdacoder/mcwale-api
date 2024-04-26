package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPgDataBase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
