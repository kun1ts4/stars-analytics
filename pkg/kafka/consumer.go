// Package kafka предоставляет клиенты для работы с Kafka.
package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Consumer потребляет сообщения из Kafka.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer создает новый Consumer.
func NewConsumer(brokers []string, topic /*,groupID*/ string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		// GroupID: groupID,
	})
	return &Consumer{reader: reader}
}

// Read читает сообщение из Kafka.
func (c *Consumer) Read(ctx context.Context) ([]byte, error) {
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading message from Kafka: %w", err)
	}
	return message.Value, nil
}

// Close закрывает Consumer.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
