package metrics

import (
	"database/sql"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// gRPC Server Metrics
	Requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"service", "method"},
	)

	Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_grpc_errors_total",
			Help: "Total number of gRPC requests resulting in error",
		},
		[]string{"service", "method", "code"},
	)

	Latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_grpc_request_duration_seconds",
			Help:    "gRPC request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)

	// DB Metrics
	DBLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_db_duration_seconds",
			Help:    "Database query latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	DBErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_db_errors_total",
			Help: "Total database errors",
		},
		[]string{"operation", "error"},
	)

	// External Services (Outgoing) Metrics
	ExternalLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_external_duration_seconds",
			Help:    "External service latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)

	ExternalErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_external_errors_total",
			Help: "Total external service errors",
		},
		[]string{"service", "method", "code"},
	)
)

// Init registers metrics with the default Prometheus registry.
func Init() {
	prometheus.MustRegister(Requests, Errors, Latency, DBLatency, DBErrors, ExternalLatency, ExternalErrors)
}

// RegisterDBStats registers SQL connection pool metrics.
func RegisterDBStats(db *sql.DB) {
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{Name: "go_db_pool_open_connections", Help: "Number of open connections"},
		func() float64 { return float64(db.Stats().OpenConnections) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{Name: "go_db_pool_in_use", Help: "Number of connections currently in use"},
		func() float64 { return float64(db.Stats().InUse) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{Name: "go_db_pool_idle", Help: "Number of idle connections"},
		func() float64 { return float64(db.Stats().Idle) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{Name: "go_db_pool_wait_count", Help: "Total number of connections waited for"},
		func() float64 { return float64(db.Stats().WaitCount) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{Name: "go_db_pool_wait_duration_seconds", Help: "Total time blocked waiting for connections"},
		func() float64 { return db.Stats().WaitDuration.Seconds() },
	))
}

// Handler returns the HTTP handler for /metrics.
func Handler() http.Handler {
	return promhttp.Handler()
}
