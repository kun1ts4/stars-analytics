// Package storage предоставляет функции для миграции и хранения данных.
package storage

import (
	"github.com/kun1ts4/stars-analytics/internal/storage/models"
	"gorm.io/gorm"
)

// Migrate выполняет миграцию базы данных.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.HourlyAggregate{},
	)
}
