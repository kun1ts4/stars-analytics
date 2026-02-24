// Package server предоставляет функционал для запуска gRPC сервера.
package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/config"
	"github.com/kun1ts4/stars-analytics/internal/prometheus"
	gormrepo "github.com/kun1ts4/stars-analytics/internal/storage/gorm"
	"github.com/kun1ts4/stars-analytics/pkg/logger"
	"github.com/kun1ts4/stars-analytics/pkg/pb/github.com/kun1ts4/stars-analytics/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// Manager управляет жизненным циклом gRPC и HTTP серверов.
type Manager struct {
	grpcServer     *grpc.Server
	metricsServer  *http.Server
	grpcListener   net.Listener
	grpcServerErr  chan error
	metricsAddress string
	grpcAddress    string
}

// NewManager создает новый Manager.
func NewManager(cfg *config.Config, db *gorm.DB) (*Manager, error) {
	prometheus.Init()

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	prometheus.RegisterDBStats(sqlDB)

	repo := gormrepo.NewStatsRepo(db)
	srv := &Server{
		UnimplementedStatsServer: &proto.UnimplementedStatsServer{},
		Repo:                     repo,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(MetricsInterceptor),
	)
	proto.RegisterStatsServer(grpcServer, srv)

	metricsServer := &http.Server{
		Addr:    ":9090",
		Handler: prometheus.Handler(),
	}

	listener, err := net.Listen("tcp", cfg.GRPC.Address())
	if err != nil {
		return nil, err
	}

	return &Manager{
		grpcServer:     grpcServer,
		metricsServer:  metricsServer,
		grpcListener:   listener,
		grpcServerErr:  make(chan error, 1),
		metricsAddress: metricsServer.Addr,
		grpcAddress:    cfg.GRPC.Address(),
	}, nil
}

// Start запускает оба сервера.
func (sm *Manager) Start(ctx context.Context) {
	go func() {
		logger.WithFields(logrus.Fields{
			"address": sm.metricsAddress,
		}).Info("Metrics server listening")
		if err := sm.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("failed to start prometheus server")
		}
	}()

	go func() {
		logger.WithFields(logrus.Fields{
			"address": sm.grpcAddress,
		}).Info("gRPC server listening")
		sm.grpcServerErr <- sm.grpcServer.Serve(sm.grpcListener)
	}()

	<-ctx.Done()
	sm.Shutdown(context.Background())
}

// Shutdown выполняет корректное завершение обоих серверов.
func (sm *Manager) Shutdown(ctx context.Context) {
	logger.Info("shutdown signal received, initiating graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownCancel()

	go func() {
		sm.grpcServer.GracefulStop()
		sm.grpcServerErr <- nil
	}()

	select {
	case err := <-sm.grpcServerErr:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.WithError(err).Error("grpc server error")
		}
		logger.Info("gRPC server stopped")
	case <-shutdownCtx.Done():
		logger.Warn("gRPC server graceful shutdown timeout, forcing stop")
		sm.grpcServer.Stop()
	}

	// Shutdown Prometheus HTTP server with timeout
	metricsCtx, metricsCancel := context.WithTimeout(ctx, 10*time.Second)
	defer metricsCancel()
	if err := sm.metricsServer.Shutdown(metricsCtx); err != nil {
		logger.WithError(err).Warn("error shutting down metrics server")
	} else {
		logger.Info("metrics server stopped")
	}
}
