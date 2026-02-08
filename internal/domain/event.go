// Package domain определяет основные бизнес-сущности.
package domain

import (
	"time"
)

// Event представляет событие звезды GitHub.
type Event struct {
	ID     string
	Action ActionType

	RepoID     int64
	RepoName   string
	ActorLogin string

	CreatedAt time.Time
}

// ActionType представляет тип действия в событии.
type ActionType string

const (
	// ActionStarred является типом действия для звезды репозитория.
	ActionStarred ActionType = "stared"
)
