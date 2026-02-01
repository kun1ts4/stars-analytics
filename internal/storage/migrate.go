package storage

import (
	"github.com/kun1ts4/stars-analytics/internal/storage/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.HourlyAggregate{},
	)
}