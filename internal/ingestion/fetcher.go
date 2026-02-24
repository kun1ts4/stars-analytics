// Package ingestion обрабатывает получение и обработку данных архива GitHub.
package ingestion

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
)

// KafkaProducer определяет интерфейс для производителя Kafka.
type KafkaProducer interface {
	Send(ctx context.Context, key string, value []byte) error
	Close() error
}

// GHArchiveFetcher получает события GitHub из GH Archive.
type GHArchiveFetcher struct {
	httpClient    *http.Client
	lastProcessed time.Time
	producer      KafkaProducer
	config        config.IngestionConfig
}

// NewGHArchiveFetcher создает новый GHArchiveFetcher.
func NewGHArchiveFetcher(
	httpClient *http.Client,
	lastProcessed time.Time,
	producer KafkaProducer,
	cfg config.IngestionConfig,
) *GHArchiveFetcher {
	return &GHArchiveFetcher{
		httpClient:    httpClient,
		lastProcessed: lastProcessed,
		producer:      producer,
		config:        cfg,
	}
}

// Run запускает цикл получения.
func (f *GHArchiveFetcher) Run(ctx context.Context) error {
	pollInterval := time.Duration(f.config.PollIntervalSec) * time.Second

	for {
		select {
		case <-ctx.Done():
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
			defer shutdownCancel()

			done := make(chan struct{})
			go func() {
				err := f.producer.Close()
				if err != nil {
					logger.WithError(err).Error("failed to close producer")
				}
				close(done)
			}()

			select {
			case <-done:
				logger.Info("finished closing producer")
			case <-shutdownCtx.Done():
				logger.Info("shutting down producer")
			}
		default:
			nextHour := f.lastProcessed.Add(time.Hour)
			if time.Since(nextHour) >= time.Hour {
				if err := f.fetchHour(nextHour); err != nil {
					logger.WithError(err).Warn("fetch failed")
				}
			} else {
				time.Sleep(pollInterval)
			}
		}
	}
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
			logger.WithError(err).Warn("failed to close body")
		}
	}()

	if err := f.processStream(body); err != nil {
		return err
	}

	f.lastProcessed = t
	logger.WithFields(logrus.Fields{
		"hour": t.Format("2006-01-02 15"),
	}).Info("finished processing hour")
	return nil
}

func (f *GHArchiveFetcher) downloadHour(date time.Time) (io.ReadCloser, error) {
	url := fmt.Sprintf(
		"%s%s-%02d.json.gz",
		f.config.GHArchiveURL,
		date.Format("2006-01-02"),
		date.Hour(),
	)

	logger.WithFields(logrus.Fields{
		"url": url,
	}).Info("downloading")

	resp, err := f.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if err := resp.Body.Close(); err != nil {
			logger.WithError(err).Warn("failed to close response body")
		}
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	return resp.Body, nil
}
