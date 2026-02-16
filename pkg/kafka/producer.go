package kafka

import (
	"context"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/segmentio/kafka-go"
)

// Producer производит сообщения в Kafka.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer создает новый Producer.
func NewProducer(brokers []string, topic string, cfg config.KafkaProducerConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,

			BatchSize:    cfg.BatchSize,
			BatchTimeout: time.Duration(cfg.BatchTimeoutMs) * time.Millisecond,
			RequiredAcks: kafka.RequiredAcks(cfg.RequiredAcks),
			MaxAttempts:  cfg.MaxAttempts,
		},
	}
}

// Send отправляет сообщение в Kafka.
func (p *Producer) Send(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: value,
			Time:  time.Now(),
		},
	)
}

// Close закрывает Producer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
