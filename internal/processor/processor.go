package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/kun1ts4/stars-analytics/internal/dto"
	"github.com/kun1ts4/stars-analytics/internal/storage"
	"github.com/kun1ts4/stars-analytics/pkg/kafka"
)

type Processor struct {
	Consumer  *kafka.Consumer
	StatsRepo *storage.StatsGormRepo
}

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

func (p *Processor) ProcessEvent(event domain.Event) error {
	return p.StatsRepo.UpdateCounts(event)
}
