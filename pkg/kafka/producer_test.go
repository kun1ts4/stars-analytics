package kafka

import (
	"testing"

	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/stretchr/testify/require"
)

func TestNewProducer(t *testing.T) {
	cfg := config.KafkaProducerConfig{
		BatchSize:      100,
		BatchTimeoutMs: 100,
		MaxAttempts:    3,
		RequiredAcks:   1,
	}

	tests := []struct {
		name    string
		brokers []string
		topic   string
	}{
		{
			name:    "single_broker",
			brokers: []string{"localhost:9092"},
			topic:   "test-topic",
		},
		{
			name:    "multiple_brokers",
			brokers: []string{"broker1:9092", "broker2:9092", "broker3:9092"},
			topic:   "events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			producer := NewProducer(tt.brokers, tt.topic, cfg)
			require.NotNil(t, producer)
			require.NotNil(t, producer.writer)
			require.Equal(t, tt.topic, producer.writer.Topic)
		})
	}
}
