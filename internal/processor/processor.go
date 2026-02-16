// Package processor обрабатывает события из Kafka и сохраняет агрегированные данные.
package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/kun1ts4/stars-analytics/internal/dto"
)

// KafkaConsumer определяет интерфейс для потребителя Kafka.
type KafkaConsumer interface {
	Read(ctx context.Context) ([]byte, error)
	Close() error
}

// Processor обрабатывает события.
type Processor struct {
	Consumer  KafkaConsumer
	StatsRepo domain.StatsRepo
}

// Run запускает обработку событий.
func (p *Processor) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := p.Consumer.Read(ctx)
			if err != nil {
				return err
			}
			kafkaEvent := dto.KafkaEvent{}
			err = json.Unmarshal(msg, &kafkaEvent)
			if err != nil {
				return fmt.Errorf("unmarshalling Kafka event: %w", err)
			}
			event := kafkaEvent.ToDomain()
			err = p.ProcessEvent(event)
			if err != nil {
				return err
			}
		}
	}
}

// ProcessEvent обрабатывает отдельное событие.
func (p *Processor) ProcessEvent(event domain.Event) error {
	return p.StatsRepo.UpdateCounts(event)
}
