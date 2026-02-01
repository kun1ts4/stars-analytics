package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/ingestion/dto"
	"github.com/kun1ts4/stars-analytics/pkg/kafka"
)

type GHArchiveFetcher struct {
	httpClient    *http.Client
	lastProcessed time.Time
	producer      *kafka.Producer
}

func NewGHArchiveFetcher(httpClient *http.Client, lastProcessed time.Time, producer *kafka.Producer) *GHArchiveFetcher {
	return &GHArchiveFetcher{
		httpClient:    httpClient,
		lastProcessed: lastProcessed,
		producer:      producer,
	}
}

func (f *GHArchiveFetcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	// Initial run immediately
	if err := f.processCurrentHour(); err != nil {
		log.Printf("initial fetch failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := f.processCurrentHour(); err != nil {
				log.Printf("fetch failed: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (f *GHArchiveFetcher) processCurrentHour() error {
	hourToFetch := time.Now().UTC().Add(-1 * time.Hour)
	if err := f.fetchHour(hourToFetch); err != nil {
		return fmt.Errorf("failed to fetch hour %s: %w", hourToFetch, err)
	}
	return nil
}

func (f *GHArchiveFetcher) fetchHour(t time.Time) error {
	if time.Since(t) < time.Hour {
		return fmt.Errorf("data not ready yet, need to wait")
	}
	body, err := f.downloadHour(t)
	if err != nil {
		return err
	}
	defer func() {
		if err := body.Close(); err != nil {
			log.Printf("failed to close body: %v", err)
		}
	}()

	if err := f.processStream(body); err != nil {
		return err
	}

	f.lastProcessed = t
	log.Printf("Finished hour: %s", t.Format("2006-01-02 15"))
	return nil
}

func (f *GHArchiveFetcher) downloadHour(date time.Time) (io.ReadCloser, error) {
	url := fmt.Sprintf("https://data.gharchive.org/%s-%d.json.gz",
		date.Format("2006-01-02"), date.Hour())

	resp, err := f.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return resp.Body, nil
	return resp.Body, nil
}

func (f *GHArchiveFetcher) processStream(gzStream io.Reader) error {
	events := make(chan dto.EventDTO, 10000)

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
					data, err := json.Marshal(event)
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
