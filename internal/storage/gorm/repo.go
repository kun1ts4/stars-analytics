// Package gorm предоставляет реализацию репозитория с использованием GORM.
package gorm

import (
	"fmt"
	"sync"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/kun1ts4/stars-analytics/internal/storage/models"
	"github.com/kun1ts4/stars-analytics/pkg/pb/github.com/kun1ts4/stars-analytics/proto"
	"gorm.io/gorm"
)

// StatsRepo реализует domain.StatsRepo с использованием GORM.
type StatsRepo struct {
	db *gorm.DB
	mu sync.Mutex
}

// NewStatsRepo создаёт новый репозиторий статистики.
func NewStatsRepo(db *gorm.DB) domain.StatsRepo {
	return &StatsRepo{db: db}
}

// UpdateCounts обновляет счетчики для события.
func (r *StatsRepo) UpdateCounts(event domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hourBucket := event.CreatedAt.Truncate(time.Hour)

	result := r.db.Model(&models.HourlyAggregate{}).
		Where("repo_id = ? AND hour = ?", event.RepoID, hourBucket).
		Update("stars", gorm.Expr("stars + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("updating repo stars: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		r.db.Create(&models.HourlyAggregate{
			RepoID:   event.RepoID,
			RepoName: event.RepoName,
			Stars:    1,
			Hour:     hourBucket,
		})
	}

	return nil
}

// GetTopN возвращает топ N репозиториев.
func (r *StatsRepo) GetTopN(count int) ([]*proto.Repo, error) {
	hourBucket := time.Now().UTC().Add(time.Hour * -1).Truncate(time.Hour)

	var aggregates []models.HourlyAggregate
	result := r.db.Where("hour = ?", hourBucket).Limit(count).Order("stars desc").Find(&aggregates)
	if result.Error != nil {
		return nil, fmt.Errorf("getting top n: %w", result.Error)
	}

	repos := make([]*proto.Repo, len(aggregates))
	for i, agg := range aggregates {
		repos[i] = &proto.Repo{
			Name:          agg.RepoName,
			StarsLastHour: uint64(agg.Stars),
			TotalStars:    10000, // TODO total
		}
	}

	return repos, nil
}
