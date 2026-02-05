package ingestion

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
