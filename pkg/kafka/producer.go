package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer производит сообщения в Kafka.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer создает новый Producer.
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,

			BatchSize:    100,
			BatchTimeout: 100 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
			MaxAttempts:  3,
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
