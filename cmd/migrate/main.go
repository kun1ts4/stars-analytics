package main

import (
	"log"

	"github.com/kun1ts4/stars-analytics/internal/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Migrate(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed")
}
