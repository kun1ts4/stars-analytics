package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/kun1ts4/stars-analytics/internal/storage/models"
	"github.com/kun1ts4/stars-analytics/pkg/pb/github.com/kun1ts4/stars-analytics/proto"
	"gorm.io/gorm"
)

// StatsGormRepo реализует StatsRepo с использованием GORM.
type StatsGormRepo struct {
	Db *gorm.DB
	Mu sync.Mutex
}

// StatsRepo определяет интерфейс для репозитория статистики.
type StatsRepo interface {
	SaveHourlyAggregate(*[]models.HourlyAggregate) error
	GetTopN(int) ([]*proto.Repo, error)
}

// UpdateCounts обновляет счетчики для события.
func (r *StatsGormRepo) UpdateCounts(event domain.Event) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	hourBucket := event.CreatedAt.Truncate(time.Hour)

	result := r.Db.Model(&models.HourlyAggregate{}).
		Where("repo_id = ? AND hour = ?", event.RepoID, hourBucket).
		Update("stars", gorm.Expr("stars + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("updating repo stars: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		r.Db.Create(&models.HourlyAggregate{
			RepoID:   event.RepoID,
			RepoName: event.RepoName,
			Stars:    1,
			Hour:     hourBucket,
		})
	}

	return nil
}

// GetTopN возвращает топ N репозиториев.
func (r *StatsGormRepo) GetTopN(count int) ([]*proto.Repo, error) {
	hourBucket := time.Now().UTC().Add(time.Hour * -1).Truncate(time.Hour)

	var aggregates []models.HourlyAggregate
	result := r.Db.Where("hour = ?", hourBucket).Limit(count).Order("stars desc").Find(&aggregates)
	if result.Error != nil {
		return nil, fmt.Errorf("getting top n: %v", result.Error)
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
