// /api-graphql/metrics.go
package graphql

import (
	"net/http"

	"github.com/goletan/observability/metrics"
	"github.com/goletan/observability/utils"
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

// Security Tool: Scrub sensitive data
var (
	scrubber = utils.NewScrubber()
)

func InitMetrics() {
	metrics.NewManager().Register(&GraphQLMetrics{})
}

func (em *GraphQLMetrics) Register() error {
	if err := prometheus.Register(RequestDuration); err != nil {
		return err
	}

	return nil
}

// ObserveRequestDuration records the duration of HTTP requests.
func ObserveRequestDuration(method, endpoint string, status int, duration float64) {
	scrubbedMethod := scrubber.Scrub(method)
	scrubbedEndpoint := scrubber.Scrub(endpoint)
	RequestDuration.WithLabelValues(scrubbedMethod, scrubbedEndpoint, http.StatusText(status)).Observe(duration)
}
