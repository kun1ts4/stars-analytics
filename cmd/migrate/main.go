package main

import (
	"log"

	"github.com/kun1ts4/stars-analytics/internal/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=postgres user=postgres password=postgres dbname=stars_analytics port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := storage.Migrate(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed")
}
