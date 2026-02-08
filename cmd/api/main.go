// cmd/api/main.go
// Команда api запускает gRPC сервер для доступа к агрегированным данным
package main

import (
	"net"
	"sync"

	server "github.com/kun1ts4/stars-analytics/internal/api"
	"github.com/kun1ts4/stars-analytics/internal/storage"
	"github.com/kun1ts4/stars-analytics/pkg/pb/github.com/kun1ts4/stars-analytics/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(
		postgres.Open(
			"host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable",
		),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}

	repo := storage.StatsGormRepo{
		Db: db,
		Mu: sync.Mutex{},
	}

	srv := server.Server{
		UnimplementedStatsServer: &proto.UnimplementedStatsServer{},
		Repo:                     &repo,
	}

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic("failed to listen")
	}

	s := grpc.NewServer()
	proto.RegisterStatsServer(s, &srv)
	if err := s.Serve(listen); err != nil {
		panic("failed to serve")
	}
}
