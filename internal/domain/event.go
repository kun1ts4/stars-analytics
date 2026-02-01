package domain

import (
	"time"
)

type Event struct {
	ID      string
	Type    string
	Payload struct {
		Action string
	}
	Repo struct {
		ID   int
		Name string
		URL  string
	}
	Actor struct {
		ID    int
		Login string
		URL   string
	}
	CreatedAt time.Time
}
