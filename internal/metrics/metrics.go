package metrics

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// IngestLatency is a histogram vector labeled by HTTP status.
var IngestLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "ingest_latency_ms",
		Help:    "Latency of /ingest endpoint in milliseconds",
		Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1ms → ~512ms
	},
	[]string{"status"},
)

// Setup registers the histogram and returns a Gin middleware.
func Setup(r *gin.Engine) gin.HandlerFunc {
	// 1. Register the metric (idempotent if called once at startup)
	prometheus.MustRegister(IngestLatency)

	// 2. Return middleware that observes request latency
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		status := fmt.Sprint(c.Writer.Status())
		elapsed := float64(time.Since(start).Milliseconds())
		IngestLatency.WithLabelValues(status).Observe(elapsed)
	}
}
