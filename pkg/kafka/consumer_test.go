package kafka

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConsumer(t *testing.T) {
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
			brokers: []string{"broker1:9092", "broker2:9092"},
			topic:   "events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := NewConsumer(tt.brokers, tt.topic)
			require.NotNil(t, consumer)
			require.NotNil(t, consumer.reader)
		})
	}
}
