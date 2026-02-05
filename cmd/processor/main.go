// cmd/processor/main.go
// Команда processor потребляет события из Kafka, агрегирует данные и сохраняет результаты в Postgres
package main

import (
	"context"
	"fmt"

	processor "github.com/kun1ts4/stars-analytics/internal/processor"
	"github.com/kun1ts4/stars-analytics/internal/storage"
	"github.com/kun1ts4/stars-analytics/pkg/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open("host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	repo := &storage.StatsGormRepo{Db: db}
	consumer := kafka.NewConsumer([]string{"kafka:9092"}, "github.events")

	proc := processor.Processor{
		Consumer:  consumer,
		StatsRepo: repo,
	}
	err = proc.Run(context.Background())
	if err != nil {
		fmt.Println(fmt.Errorf("running processor: %v", err))
	}
}
