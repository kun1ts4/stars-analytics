package ingestion

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/kun1ts4/stars-analytics/internal/dto"
)

func (f *GHArchiveFetcher) processStream(gzStream io.Reader) error {
	events := make(chan dto.GHEvent, 10000)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for event := range events {
				if event.Type == "WatchEvent" && event.Payload.Action == "started" {
					if err := event.Validate(); err != nil {
						log.Printf("invalid event: %v", err)
						continue
					}
					kafkaMessage, _ := dto.ToKafkaEvent(event)
					data, err := json.Marshal(kafkaMessage)
					if err != nil {
						log.Printf("failed to marshal event: %v", err)
						continue
					}
					if err := f.producer.Send(context.Background(), event.ID, data); err != nil {
						log.Printf("failed to send event to kafka: %v", err)
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
