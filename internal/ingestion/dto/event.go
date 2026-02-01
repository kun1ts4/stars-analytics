package dto

import (
	"errors"
	"time"
)

type EventDTO struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Payload struct {
		Action string `json:"action"`
	} `json:"payload"`
	Repo struct {
		ID   int    `json:"id"`
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

func (e EventDTO) Validate() error {
	if e.ID == "" {
		return errors.New("ID is required")
	}
	if e.Type == "" {
		return errors.New("Type is required")
	}
	if e.Repo.Name == "" {
		return errors.New("Repo Name is required")
	}
	if e.Actor.Login == "" {
		return errors.New("Actor Login is required")
	}
	return nil
}
