// Package server предоставляет обработчики gRPC сервера для сервиса статистики.
package server

import (
	"context"

	"github.com/kun1ts4/stars-analytics/internal/storage"
	"github.com/kun1ts4/stars-analytics/pkg/pb/github.com/kun1ts4/stars-analytics/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server реализует интерфейс StatsServer.
type Server struct {
	*proto.UnimplementedStatsServer
	Repo *storage.StatsGormRepo
}

// TopN возвращает топ N репозиториев по звездам.
func (s *Server) TopN(_ context.Context, req *proto.NRequest) (*proto.TopResponse, error) {
	repos, err := s.Repo.GetTopN(int(req.N))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.TopResponse{Repos: repos}, nil
}
