package domain

import (
	"time"
)

type Event struct {
	ID     string
	Action ActionType

	RepoID     int64
	RepoName   string
	ActorLogin string

	CreatedAt time.Time
}

type ActionType string

const (
	ActionStarred ActionType = "stared"
)
