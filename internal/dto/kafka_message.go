package dto

import (
	"fmt"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/domain"
)

// KafkaEvent представляет сообщение события для Kafka.
type KafkaEvent struct {
	EventID   string            `json:"event_id"`
	Action    domain.ActionType `json:"action"`
	RepoID    int64             `json:"repo_id"`
	RepoName  string            `json:"repo_name"`
	UserLogin string            `json:"user_login"`
	Timestamp time.Time         `json:"timestamp"`
}

// ToDomain преобразует KafkaEvent в доменное Event.
func (e KafkaEvent) ToDomain() domain.Event {
	event := domain.Event{
		ID:         e.EventID,
		Action:     e.Action,
		RepoID:     e.RepoID,
		RepoName:   e.RepoName,
		ActorLogin: e.UserLogin,
		CreatedAt:  e.Timestamp,
	}

	return event
}

// ToKafkaEvent преобразует GHEvent в KafkaEvent.
func ToKafkaEvent(gh GHEvent) (*KafkaEvent, error) {
	switch gh.Type {
	case "WatchEvent":
		if gh.Payload.Action != "started" {
			return nil, fmt.Errorf("unsupported action for WatchEvent: %s", gh.Payload.Action)
		}
		return &KafkaEvent{
			EventID:   gh.ID,
			Action:    domain.ActionStarred,
			RepoID:    gh.Repo.ID,
			RepoName:  gh.Repo.Name,
			UserLogin: gh.Actor.Login,
			Timestamp: gh.CreatedAt,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported event type: %s", gh.Type)
	}
}
