// cmd/ingestion/main.go
// Команда ingestion собирает события из GitHub API и отправляет их в Kafka
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/ingestion"
	"github.com/kun1ts4/stars-analytics/pkg/kafka"
)

func main() {
	httpClient := &http.Client{}
	lastProceed := time.Now().UTC().Add(-2 * time.Hour)

	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "kafka:9092"
	}
	brokers := strings.Split(brokersStr, ",")

	topic := os.Getenv("TOPIC")
	if topic == "" {
		topic = "github.events"
	}

	producer := kafka.NewProducer(brokers, topic)

	fetcher := ingestion.NewGHArchiveFetcher(httpClient, lastProceed, producer)
	if err := fetcher.Run(context.Background()); err != nil {
		log.Fatalf("Failed to run fetcher: %v", err)
	}
}
