// Package models определяет модели данных для хранения.
package models

import "time"

// HourlyAggregate представляет агрегированные данные за час.
type HourlyAggregate struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	RepoID   int64  `gorm:"not null;uniqueIndex:idx_repo_hour"`
	RepoName string `gorm:"type:varchar(255);not null"`
	Stars    int    `gorm:"default:0"`
	// Forks    int       `gorm:"default:0"`
	Hour time.Time `gorm:"not null;uniqueIndex:idx_repo_hour;index:,sort:desc"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
