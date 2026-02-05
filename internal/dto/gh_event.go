package dto

import (
	"errors"
	"time"
)

type GHEvent struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload struct {
		Action string `json:"action"`
	} `json:"payload"`
	Repo struct {
		ID   int64    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	Actor struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		URL   string `json:"url"`
	}
	CreatedAt time.Time `json:"created_at"`
}

func (e GHEvent) Validate() error {
	if e.ID == "" {
		return errors.New("ID is required")
	}
	if e.Type == "" {
		return errors.New("type is required")
	}
	if e.Repo.Name == "" {
		return errors.New("repo Name is required")
	}
	if e.Actor.Login == "" {
		return errors.New("actor Login is required")
	}
	return nil
}
