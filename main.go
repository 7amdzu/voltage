package main

import (
	_ "embed"

	"github.com/7amdzu/voltage/internal/handler"
	"github.com/7amdzu/voltage/internal/metrics"
	"github.com/gin-gonic/gin"
	ginmetrics "github.com/penglongli/gin-metrics/ginmetrics"
)

//go:embed sbom.spdx.json
var sbom []byte

func main() {
	r := gin.New()

	// A) Built-in Prometheus /metrics from gin-metrics
	monitor := ginmetrics.GetMonitor()
	monitor.SetMetricPath("/metrics")
	monitor.Use(r)

	// B) Custom histogram for /ingest latency
	r.Use(metrics.Setup(r))

	// Endpoints
	r.POST("/ingest", handler.Ingest)
	r.GET("/sbom", func(c *gin.Context) {
		c.Data(200, "application/json", sbom)
	})
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
