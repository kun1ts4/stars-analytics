package gorm

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kun1ts4/stars-analytics/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestNewStatsRepo(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := NewStatsRepo(db)
	assert.NotNil(t, repo)
}

func TestUpdateCounts_NewRecord(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewStatsRepo(db)

	event := domain.Event{
		ID:         "123",
		Action:     domain.ActionStarred,
		RepoID:     1,
		RepoName:   "test/repo",
		ActorLogin: "user",
		CreatedAt:  time.Date(2024, 1, 1, 15, 30, 0, 0, time.UTC),
	}

	hourBucket := event.CreatedAt.Truncate(time.Hour)

	// Ожидаем UPDATE который не найдет записей (GORM добавляет updated_at)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "hourly_aggregates"`).
		WithArgs(1, sqlmock.AnyArg(), 1, hourBucket).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// Ожидаем INSERT новой записи
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "hourly_aggregates"`).
		WithArgs(event.RepoID, event.RepoName, 1, hourBucket, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.UpdateCounts(event)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCounts_ExistingRecord(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewStatsRepo(db)

	hourBucket := time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC)

	event := domain.Event{
		ID:         "456",
		Action:     domain.ActionStarred,
		RepoID:     1,
		RepoName:   "test/repo",
		ActorLogin: "user2",
		CreatedAt:  hourBucket.Add(30 * time.Minute),
	}

	// Ожидаем UPDATE который обновит 1 запись (GORM добавляет updated_at)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "hourly_aggregates"`).
		WithArgs(1, sqlmock.AnyArg(), 1, hourBucket).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.UpdateCounts(event)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCounts_DatabaseError(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewStatsRepo(db)

	event := domain.Event{
		ID:         "123",
		RepoID:     1,
		RepoName:   "test/repo",
		ActorLogin: "user",
		CreatedAt:  time.Now(),
	}

	hourBucket := event.CreatedAt.Truncate(time.Hour)

	// Симулируем ошибку базы данных (GORM добавляет updated_at)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "hourly_aggregates"`).
		WithArgs(1, sqlmock.AnyArg(), 1, hourBucket).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()

	err := repo.UpdateCounts(event)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopN(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewStatsRepo(db)

	hourBucket := time.Now().UTC().Add(time.Hour * -1).Truncate(time.Hour)

	rows := sqlmock.NewRows([]string{"id", "repo_id", "repo_name", "stars", "hour", "created_at", "updated_at"}).
		AddRow(1, 2, "repo2/test", 200, hourBucket, time.Now(), time.Now()).
		AddRow(2, 4, "repo4/test", 150, hourBucket, time.Now(), time.Now()).
		AddRow(3, 1, "repo1/test", 100, hourBucket, time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM "hourly_aggregates"`).
		WithArgs(hourBucket, 3).
		WillReturnRows(rows)

	repos, err := repo.GetTopN(3)
	require.NoError(t, err)
	require.Len(t, repos, 3)

	// Проверяем данные
	assert.Equal(t, "repo2/test", repos[0].Name)
	assert.Equal(t, uint64(200), repos[0].StarsLastHour)

	assert.Equal(t, "repo4/test", repos[1].Name)
	assert.Equal(t, uint64(150), repos[1].StarsLastHour)

	assert.Equal(t, "repo1/test", repos[2].Name)
	assert.Equal(t, uint64(100), repos[2].StarsLastHour)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopN_EmptyResult(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewStatsRepo(db)

	hourBucket := time.Now().UTC().Add(time.Hour * -1).Truncate(time.Hour)

	rows := sqlmock.NewRows([]string{"id", "repo_id", "repo_name", "stars", "hour", "created_at", "updated_at"})

	mock.ExpectQuery(`SELECT \* FROM "hourly_aggregates"`).
		WithArgs(hourBucket, 10).
		WillReturnRows(rows)

	repos, err := repo.GetTopN(10)
	require.NoError(t, err)
	assert.Empty(t, repos)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopN_DatabaseError(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewStatsRepo(db)

	hourBucket := time.Now().UTC().Add(time.Hour * -1).Truncate(time.Hour)

	mock.ExpectQuery(`SELECT \* FROM "hourly_aggregates"`).
		WithArgs(hourBucket, 10).
		WillReturnError(gorm.ErrInvalidDB)

	repos, err := repo.GetTopN(10)
	assert.Error(t, err)
	assert.Nil(t, repos)

	assert.NoError(t, mock.ExpectationsWereMet())
}
