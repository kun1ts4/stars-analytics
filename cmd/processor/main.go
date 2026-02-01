// cmd/processor/main.go
// Команда processor потребляет события из Kafka, агрегирует данные и сохраняет результаты в Postgres
package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStatsRepo struct {
	db *gorm.DB
}

type Example struct {
	Id   uint64 `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

func NewGormStatsRepo(db *gorm.DB) *GormStatsRepo {
	return &GormStatsRepo{db: db}
}

func (r *GormStatsRepo) ExampleMigrate() error {
	return r.db.AutoMigrate(&Example{})
}

func (r *GormStatsRepo) ExampleRead() Example {
	res := Example{}
	r.db.First(&res, "id = 2")
	return res
}

func (r *GormStatsRepo) ExampleCreate(e []Example) error {
	return r.db.Create(&e).Error
}

func main() {
	db, err := gorm.Open(postgres.Open("host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	repo := NewGormStatsRepo(db)
	err = repo.ExampleMigrate()
	if err != nil {
		panic(err)
	}
	err = repo.ExampleCreate([]Example{
		Example{
			Id:   1,
			Name: "John Doe",
		},
		Example{
			Id:   2,
			Name: "Jane Doe",
		},
	})
	if err != nil {
		panic(err)
	}

	first := repo.ExampleRead()
	fmt.Println(first)
}
