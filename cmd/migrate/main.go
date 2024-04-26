package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/razdacoder/mcwale-api/db"
	"github.com/razdacoder/mcwale-api/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}
	db, err := db.NewPgDataBase(os.Getenv("DSN"))

	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.User{})
	log.Println("Migration Complete")

}
