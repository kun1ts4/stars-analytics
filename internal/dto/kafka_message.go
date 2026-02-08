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
func ToKafkaEvent(gh GHEvent) (KafkaEvent, error) {
	event := KafkaEvent{
		EventID:   gh.ID,
		RepoID:    gh.Repo.ID,
		RepoName:  gh.Repo.Name,
		UserLogin: gh.Actor.Login,
		Timestamp: gh.CreatedAt,
	}
	switch gh.Type {
	case "WatchEvent":
		if gh.Payload.Action == "started" {
			event.Action = domain.ActionStarred
		}
	default:
		return KafkaEvent{}, fmt.Errorf("unsupported event type: %s", gh.Type)
	}
	return event, nil
}
