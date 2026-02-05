package storage

import (
	"fmt"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/kun1ts4/stars-analytics/internal/storage/models"
	"gorm.io/gorm"
)

type StatsGormRepo struct {
	Db *gorm.DB
}

type StatsRepo interface {
	SaveHourlyAggregate(*[]models.HourlyAggregate) error
}

func (r *StatsGormRepo) UpdateCounts(event domain.Event) error {
	hourBucket := event.CreatedAt.Truncate(time.Hour)

	result := r.Db.Model(&models.HourlyAggregate{}).Where("repo_id = ? AND hour = ?", event.RepoID, hourBucket).
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
