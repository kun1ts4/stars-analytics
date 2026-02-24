// cmd/processor/main.go
// Команда processor потребляет события из Kafka, агрегирует данные и сохраняет результаты в Postgres
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/kun1ts4/stars-analytics/internal/config"
	processor "github.com/kun1ts4/stars-analytics/internal/processor"
	gormrepo "github.com/kun1ts4/stars-analytics/internal/storage/gorm"
	"github.com/kun1ts4/stars-analytics/pkg/kafka"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
	"github.com/sirupsen/logrus"
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

	repo := gormrepo.NewStatsRepo(db)
	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic)

	proc := processor.Processor{
		Consumer:  consumer,
		StatsRepo: repo,
	}

	logger.WithFields(logrus.Fields{
		"topic": cfg.Kafka.Topic,
	}).Info("starting processor")
	err = proc.Run(ctx)
	if err != nil {
		logger.WithError(err).Fatal("error running processor")
	}
}
