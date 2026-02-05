package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic  /*,groupID*/ string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		//GroupID: groupID,
	})
	return &Consumer{reader: reader}
}

func (c *Consumer) Read(ctx context.Context) ([]byte, error) {
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading message from Kafka: %w", err)
	}
	return message.Value, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
