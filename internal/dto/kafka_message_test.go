package dto

import (
	"testing"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestToDomain(t *testing.T) {
	tests := []struct {
		name  string
		input KafkaEvent
		want  domain.Event
	}{
		{
			name: "basic_starred_event",
			input: KafkaEvent{
				EventID:   "3488012293",
				Action:    domain.ActionStarred,
				RepoID:    48878218,
				RepoName:  "alexylem/project",
				UserLogin: "alexylem",
				Timestamp: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
			want: domain.Event{
				ID:         "3488012293",
				Action:     domain.ActionStarred,
				RepoID:     48878218,
				RepoName:   "alexylem/project",
				ActorLogin: "alexylem",
				CreatedAt:  time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
			},
		},
		{
			name: "empty_user_login",
			input: KafkaEvent{
				EventID:   "123",
				Action:    domain.ActionStarred,
				RepoID:    1,
				RepoName:  "test/repo",
				UserLogin: "",
				Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: domain.Event{
				ID:         "123",
				Action:     domain.ActionStarred,
				RepoID:     1,
				RepoName:   "test/repo",
				ActorLogin: "",
				CreatedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToDomain()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestToKafkaEvent(t *testing.T) {
	cases := []struct {
		name    string
		input   GHEvent
		want    *KafkaEvent
		wantErr bool
	}{
		{
			name: "watch_event_started",
			input: GHEvent{
				ID:   "123",
				Type: "WatchEvent",
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   1,
					Name: "test/repo",
					URL:  "https://api.github.com/repos/test/repo",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					Login: "user",
				},
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: &KafkaEvent{
				EventID:   "123",
				Action:    domain.ActionStarred,
				RepoID:    1,
				RepoName:  "test/repo",
				UserLogin: "user",
				Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "watch_event_with_large_repo_id",
			input: GHEvent{
				ID:   "999999",
				Type: "WatchEvent",
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   999999999,
					Name: "org/large-project",
					URL:  "https://api.github.com/repos/org/large-project",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					Login: "developer",
				},
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "started",
				},
				CreatedAt: time.Date(2025, 6, 15, 12, 30, 45, 0, time.UTC),
			},
			want: &KafkaEvent{
				EventID:   "999999",
				Action:    domain.ActionStarred,
				RepoID:    999999999,
				RepoName:  "org/large-project",
				UserLogin: "developer",
				Timestamp: time.Date(2025, 6, 15, 12, 30, 45, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "unsupported_event_push",
			input: GHEvent{
				Type: "PushEvent",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unsupported_event_fork",
			input: GHEvent{
				Type: "ForkEvent",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "watch_event_non_started_action",
			input: GHEvent{
				ID:   "456",
				Type: "WatchEvent",
				Repo: struct {
					ID   int64  `json:"id"`
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					ID:   2,
					Name: "test/repo2",
				},
				Actor: struct {
					ID    int    `json:"id"`
					Login string `json:"login"`
					URL   string `json:"url"`
				}{
					Login: "user2",
				},
				Payload: struct {
					Action string `json:"action"`
				}{
					Action: "stopped",
				},
				CreatedAt: time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ToKafkaEvent(c.input)
			if c.wantErr {
				require.Error(t, err)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, c.want, got)
		})
	}
}
