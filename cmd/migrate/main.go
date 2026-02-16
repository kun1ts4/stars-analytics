// Package main предоставляет команду миграции базы данных для приложения stars-analytics.
package main

import (
	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/kun1ts4/stars-analytics/internal/storage"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.PostgresDSN()), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Fatal("failed to connect to database")
	}

	if err := storage.Migrate(db); err != nil {
		logger.WithError(err).Fatal("migration failed")
	}

	logger.Info("migration completed successfully")
}
