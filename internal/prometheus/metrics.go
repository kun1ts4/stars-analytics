// Package prometheus provides Prometheus prometheus for gRPC, database and external services.
package prometheus

import (
	"database/sql"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Requests is the total number of gRPC requests.
var Requests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_grpc_requests_total",
		Help: "Total number of gRPC requests",
	},
	[]string{"service", "method"},
)

// Errors is the total number of gRPC errors.
var Errors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_grpc_errors_total",
		Help: "Total number of gRPC requests resulting in error",
	},
	[]string{"service", "method", "code"},
)

// Latency is the gRPC request duration histogram.
var Latency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "app_grpc_request_duration_seconds",
		Help:    "gRPC request latency in seconds",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"service", "method"},
)

// DBLatency is the database query duration histogram.
var DBLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "app_db_duration_seconds",
		Help:    "Database query latency",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"operation"},
)

// DBErrors is the total number of database errors.
var DBErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_db_errors_total",
		Help: "Total database errors",
	},
	[]string{"operation", "error"},
)

// ExternalLatency is the external service call duration histogram.
var ExternalLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "app_external_duration_seconds",
		Help:    "External service latency",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"service", "method"},
)

// ExternalErrors is the total number of external service errors.
var ExternalErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_external_errors_total",
		Help: "Total external service errors",
	},
	[]string{"service", "method", "code"},
)

// Init registers prometheus with the default Prometheus registry.
func Init() {
	prometheus.MustRegister(
		Requests,
		Errors,
		Latency,
		DBLatency,
		DBErrors,
		ExternalLatency,
		ExternalErrors,
	)
}

// RegisterDBStats registers SQL connection pool prometheus.
func RegisterDBStats(db *sql.DB) {
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "go_db_pool_open_connections",
			Help: "Number of open connections",
		},
		func() float64 { return float64(db.Stats().OpenConnections) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "go_db_pool_in_use",
			Help: "Number of connections currently in use",
		},
		func() float64 { return float64(db.Stats().InUse) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{Name: "go_db_pool_idle", Help: "Number of idle connections"},
		func() float64 { return float64(db.Stats().Idle) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "go_db_pool_wait_count",
			Help: "Total number of connections waited for",
		},
		func() float64 { return float64(db.Stats().WaitCount) },
	))
	prometheus.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "go_db_pool_wait_duration_seconds",
			Help: "Total time blocked waiting for connections",
		},
		func() float64 { return db.Stats().WaitDuration.Seconds() },
	))
}

// Handler returns the HTTP handler for /prometheus.
func Handler() http.Handler {
	return promhttp.Handler()
}
