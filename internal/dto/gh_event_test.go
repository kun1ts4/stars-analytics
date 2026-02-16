package dto

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	in := GHEvent{
		ID:   "3488012293",
		Type: "WatchEvent",
		Payload: struct {
			Action string `json:"action"`
		}{
			Action: "started",
		},
		Repo: struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			ID:   48878218,
			Name: "alexylem/projectpage",
			URL:  "https://api.github.com/repos/alexylem/projectpage",
		},
		Actor: struct {
			ID    int    `json:"id"`
			Login string `json:"login"`
			URL   string `json:"url"`
		}{
			ID:    11017174,
			Login: "alexylem",
			URL:   "https://api.github.com/users/alexylem",
		},
		CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
	}
	err := in.Validate()
	require.NoError(t, err)
}

func TestValidateError(t *testing.T) {
	cases := []struct {
		name   string
		event  GHEvent
		expErr error
	}{
		{
			name: "missing ID",
			event: GHEvent{
				Type: "WatchEvent",
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   48878218,
					Name: "alexylem/projectpage",
					URL:  "https://api.github.com/repos/alexylem/projectpage",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					ID:    11017174,
					Login: "alexylem",
					URL:   "https://api.github.com/users/alexylem",
				},
				CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
			expErr: errors.New("ID is required"),
		},
		{
			name: "missing Type",
			event: GHEvent{
				ID: "3488012293",
				// Type missing
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   48878218,
					Name: "alexylem/projectpage",
					URL:  "https://api.github.com/repos/alexylem/projectpage",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					ID:    11017174,
					Login: "alexylem",
					URL:   "https://api.github.com/users/alexylem",
				},
				CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
			expErr: errors.New("type is required"),
		},
		{
			name: "missing Repo Name",
			event: GHEvent{
				ID:   "3488012293",
				Type: "WatchEvent",
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID: 48878218,
					// Name missing
					URL: "https://api.github.com/repos/alexylem/projectpage",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					ID:    11017174,
					Login: "alexylem",
					URL:   "https://api.github.com/users/alexylem",
				},
				CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
			expErr: errors.New("repo Name is required"),
		},
		{
			name: "missing Actor Login",
			event: GHEvent{
				ID:   "3488012293",
				Type: "WatchEvent",
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   48878218,
					Name: "alexylem/projectpage",
					URL:  "https://api.github.com/repos/alexylem/projectpage",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					ID:  11017174,
					URL: "https://api.github.com/users/alexylem",
				},
				CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
			expErr: errors.New("actor Login is required"),
		},
		{
			name: "missing Repo ID",
			event: GHEvent{
				ID:   "3488012293",
				Type: "WatchEvent",
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					// ID is 0
					Name: "alexylem/projectpage",
					URL:  "https://api.github.com/repos/alexylem/projectpage",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					ID:    11017174,
					Login: "alexylem",
					URL:   "https://api.github.com/users/alexylem",
				},
				CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
			expErr: errors.New("repo ID is required"),
		},
		{
			name: "missing Actor ID",
			event: GHEvent{
				ID:   "3488012293",
				Type: "WatchEvent",
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   48878218,
					Name: "alexylem/projectpage",
					URL:  "https://api.github.com/repos/alexylem/projectpage",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					// ID is 0
					Login: "alexylem",
					URL:   "https://api.github.com/users/alexylem",
				},
				CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
			expErr: errors.New("actor ID is required"),
		},
		{
			name: "missing CreatedAt",
			event: GHEvent{
				ID:   "3488012293",
				Type: "WatchEvent",
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   48878218,
					Name: "alexylem/projectpage",
					URL:  "https://api.github.com/repos/alexylem/projectpage",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					ID:    11017174,
					Login: "alexylem",
					URL:   "https://api.github.com/users/alexylem",
				},
				// CreatedAt is zero
			},
			expErr: errors.New("created_at is required"),
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			err := cs.event.Validate()
			require.EqualError(t, err, cs.expErr.Error())
		})
	}
}
