package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/razdacoder/mcwale-api/cmd/api"
	"github.com/razdacoder/mcwale-api/db"
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
	log.Println("DB Connected successfully")
	server := api.NewAPISever(":8000", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
