// cmd/ingestion/main.go
// Команда ingestion собирает события из GitHub API и отправляет их в Kafka
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/kun1ts4/stars-analytics/internal/ingestion"
	"github.com/kun1ts4/stars-analytics/pkg/kafka"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	httpClient := &http.Client{}
	lastProceed := time.Now().UTC().Add(-time.Duration(cfg.Ingestion.LookbackHours) * time.Hour)

	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Producer)

	fetcher := ingestion.NewGHArchiveFetcher(httpClient, lastProceed, producer, cfg.Ingestion)

	logger.WithFields(logrus.Fields{
		"gharchive_url":  cfg.Ingestion.GHArchiveURL,
		"lookback_hours": cfg.Ingestion.LookbackHours,
	}).Info("starting ingestion service")

	if err := fetcher.Run(context.Background()); err != nil {
		logger.WithError(err).Fatal("failed to run fetcher")
	}
}
