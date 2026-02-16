// cmd/api/main.go
// Команда api запускает gRPC сервер для доступа к агрегированным данным
package main

import (
	"net"
	"net/http"

	server "github.com/kun1ts4/stars-analytics/internal/api"
	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/kun1ts4/stars-analytics/internal/prometheus"
	gormrepo "github.com/kun1ts4/stars-analytics/internal/storage/gorm"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
	"github.com/kun1ts4/stars-analytics/pkg/pb/github.com/kun1ts4/stars-analytics/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	// Initialize Prometheus prometheus
	prometheus.Init()

	db, err := gorm.Open(
		postgres.Open(cfg.Database.DSN()),
		&gorm.Config{},
	)
	if err != nil {
		logger.WithError(err).Fatal("failed to connect database")
	}

	// Register DB connection pool stats
	sqlDB, err := db.DB()
	if err != nil {
		logger.WithError(err).Fatal("failed to get underlying sql.DB")
	}
	prometheus.RegisterDBStats(sqlDB)

	repo := gormrepo.NewStatsRepo(db)

	srv := server.Server{
		UnimplementedStatsServer: &proto.UnimplementedStatsServer{},
		Repo:                     repo,
	}

	// Start prometheus HTTP server
	go func() {
		http.Handle("/metrics", prometheus.Handler())
		metricsAddr := ":9090"
		logger.WithFields(logrus.Fields{
			"address": metricsAddr,
		}).Info("Metrics server listening")
		if err := http.ListenAndServe(metricsAddr, nil); err != nil {
			logger.WithError(err).Fatal("failed to start prometheus server")
		}
	}()

	listen, err := net.Listen("tcp", cfg.GRPC.Address())
	if err != nil {
		logger.WithError(err).Fatal("failed to listen")
	}

	// Create gRPC server with prometheus interceptor
	s := grpc.NewServer(
		grpc.UnaryInterceptor(server.MetricsInterceptor),
	)
	proto.RegisterStatsServer(s, &srv)

	logger.WithFields(logrus.Fields{
		"address": cfg.GRPC.Address(),
	}).Info("gRPC server listening")
	if err := s.Serve(listen); err != nil {
		logger.WithError(err).Fatal("failed to serve")
	}
}
