package storage

import (
	"github.com/kun1ts4/stars-analytics/internal/storage/models"
	"gorm.io/gorm"
)

type StatsGormRepo struct {
	Db *gorm.DB
}

type StatsRepo interface {
	SaveHourlyAggregate(*[]models.HourlyAggregate) error
}
