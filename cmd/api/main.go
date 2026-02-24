// cmd/api/main.go
// Команда api запускает gRPC сервер для доступа к агрегированным данным
package main

import (
	"context"
	"os/signal"
	"syscall"

	apiserver "github.com/kun1ts4/stars-analytics/internal/api"
	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	db, err := gorm.Open(
		postgres.Open(cfg.Database.DSN()),
		&gorm.Config{},
	)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect database")
	}

	serverManager, err := apiserver.NewServerManager(cfg, db)
	if err != nil {
		logger.WithError(err).Fatal("failed to create server manager")
	}

	serverManager.Start(ctx)
}
