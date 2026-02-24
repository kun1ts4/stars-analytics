// Package processor обрабатывает события из Kafka и сохраняет агрегированные данные.
package processor

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/kun1ts4/stars-analytics/internal/dto"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
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
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			done := make(chan struct{})
			go func() {
				if err := p.Consumer.Close(); err != nil {
					logger.WithError(err).Error("error closing consumer")
				}
				close(done)
			}()

			select {
			case <-done:
				logger.Info("consumer closed successfully")
			case <-shutdownCtx.Done():
				logger.Warn("consumer close timeout exceeded")
			}

			return ctx.Err()
		default:
			msg, err := p.Consumer.Read(ctx)
			if err != nil {
				logger.WithError(err).Error("error reading message")
				continue
			}
			kafkaEvent := dto.KafkaEvent{}
			err = json.Unmarshal(msg, &kafkaEvent)
			if err != nil {
				logger.WithError(err).Error("error unmarshalling message")
				continue
			}
			event := kafkaEvent.ToDomain()
			err = p.ProcessEvent(event)
			if err != nil {
				logger.WithError(err).Error("error processing event")
				continue
			}
		}
	}
}

// ProcessEvent обрабатывает отдельное событие.
func (p *Processor) ProcessEvent(event domain.Event) error {
	return p.StatsRepo.UpdateCounts(event)
}
