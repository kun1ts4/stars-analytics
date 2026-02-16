package domain

import (
	"github.com/kun1ts4/stars-analytics/pkg/pb/github.com/kun1ts4/stars-analytics/proto"
)

// StatsRepo определяет интерфейс для репозитория статистики.
type StatsRepo interface {
	UpdateCounts(event Event) error
	GetTopN(count int) ([]*proto.Repo, error)
}
