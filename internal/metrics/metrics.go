// /api-graphql/internal/metrics/metrics.go
package metrics

import (
	"net/http"

	observability "github.com/goletan/observability/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type GraphQLMetrics struct{}

// HTTP Metrics: Track HTTP requests and errors.
var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "goletan",
			Subsystem: "graphql",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
		},
		[]string{"method", "endpoint", "status"},
	)
)

// InitMetrics initializes and registers GraphQL metrics with the observability.
func InitMetrics(obs *observability.Observability) *GraphQLMetrics {
	met := &GraphQLMetrics{}
	obs.Metrics.Register(met)
	return met
}

// Register registers the GraphQL metrics with Prometheus.
func (em *GraphQLMetrics) Register() error {
	if err := prometheus.Register(RequestDuration); err != nil {
		return err
	}
	return nil
}

// ObserveRequestDuration records the duration of HTTP requests.
func ObserveRequestDuration(method, endpoint string, status int, duration float64) {
	RequestDuration.WithLabelValues(method, endpoint, http.StatusText(status)).Observe(duration)
}
