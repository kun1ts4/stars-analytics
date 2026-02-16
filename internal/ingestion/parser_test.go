package ingestion

import (
	"testing"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/dto"
	"github.com/stretchr/testify/require"
)

func TestParseEvent(t *testing.T) {
	in := []byte("{\"id\":\"3488012293\",\"type\":\"WatchEvent\",\"actor\":{\"id\":11017174,\"login\":\"alexylem\",\"gravatar_id\":\"\",\"url\":\"https://api.github.com/users/alexylem\",\"avatar_url\":\"https://avatars.githubusercontent.com/u/11017174?\"},\"repo\":{\"id\":48878218,\"name\":\"alexylem/projectpage\",\"url\":\"https://api.github.com/repos/alexylem/projectpage\"},\"payload\":{\"action\":\"started\"},\"public\":true,\"created_at\":\"2016-01-02T15:00:03Z\"}\n")

	exp := dto.GHEvent{
		ID:   "3488012293",
		Type: "WatchEvent",
		Actor: struct {
			ID    int    `json:"id"`
			Login string `json:"login"`
			URL   string `json:"url"`
		}{
			ID:    11017174,
			Login: "alexylem",
			URL:   "https://api.github.com/users/alexylem",
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
		Payload: struct {
			Action string `json:"action"`
		}{
			Action: "started",
		},
		CreatedAt: time.Date(2016, 1, 2, 15, 0, 3, 0, time.UTC),
	}

	got, err := ParseEvent(in)
	require.NoError(t, err)
	require.Equal(t, exp, got)
}
