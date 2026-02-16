package ingestion

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/kun1ts4/stars-analytics/internal/dto"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
)

func (f *GHArchiveFetcher) processStream(gzStream io.Reader) error {
	events := make(chan dto.GHEvent, f.config.ChannelSize)

	var wg sync.WaitGroup
	for i := 0; i < f.config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range events {
				if event.Type == "WatchEvent" && event.Payload.Action == "started" {
					if err := event.Validate(); err != nil {
						logger.WithError(err).WithFields(logrus.Fields{
							"event_id": event.ID,
						}).Warn("invalid event")
						continue
					}
					kafkaMessage, err := dto.ToKafkaEvent(event)
					if err != nil {
						logger.WithError(err).WithFields(logrus.Fields{
							"event_id": event.ID,
						}).Warn("failed to convert event to kafka message")
						continue
					}
					data, err := json.Marshal(kafkaMessage)
					if err != nil {
						logger.WithError(err).WithFields(logrus.Fields{
							"event_id": event.ID,
						}).Warn("failed to marshal event")
						continue
					}
					if err := f.producer.Send(context.Background(), event.ID, data); err != nil {
						logger.WithError(err).WithFields(logrus.Fields{
							"event_id": event.ID,
						}).Warn("failed to send event to kafka")
					}
				}
			}
		}()
	}

	err := ParseStream(gzStream, events)
	close(events)
	wg.Wait()
	return err
}
