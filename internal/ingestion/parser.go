package ingestion

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/kun1ts4/stars-analytics/internal/dto"
)

// ParseStream reads a gzipped JSON stream and sends events to the channel
func ParseStream(r io.Reader, events chan<- dto.GHEvent) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer func() {
		if err := gz.Close(); err != nil {
			log.Printf("failed to close gzip reader: %v", err)
		}
	}()

	scanner := bufio.NewScanner(gz)
	const maxCapacity = 50 * 1024 * 1024 // 50MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	count := 0
	for scanner.Scan() {
		event, err := parseEvent(scanner.Bytes())
		if err != nil {
			log.Printf("failed to parse event: %v", err)
			continue
		}
		events <- event
		count++
	}

	if count == 0 {
		log.Printf("No events found in the stream")
	}

	return scanner.Err()
}

func parseEvent(bytes []byte) (dto.GHEvent, error) {
	event := dto.GHEvent{}
	err := json.Unmarshal(bytes, &event)
	if err != nil {
		return event, fmt.Errorf("unmarshal event: %w", err)
	}

	return event, nil
}
