package server

import (
	"context"
	"strings"
	"time"

	"github.com/kun1ts4/stars-analytics/internal/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// MetricsInterceptor собирает метрики Prometheus для gRPC запросов.
func MetricsInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	service, method := splitMethod(info.FullMethod)

	prometheus.Requests.WithLabelValues(service, method).Inc()
	prometheus.Latency.WithLabelValues(service, method).
		Observe(time.Since(start).Seconds())

	if err != nil {
		st, _ := status.FromError(err)
		prometheus.Errors.WithLabelValues(service, method, st.Code().String()).Inc()
	}

	return resp, err
}

// splitMethod extracts service and method from info.FullMethod.
// FullMethod format: "/package.Service/Method"
func splitMethod(fullMethod string) (service, method string) {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 3 {
		service = parts[1] // package.Service
		method = parts[2]
		return service, method
	}
	return "unknown", fullMethod
}
