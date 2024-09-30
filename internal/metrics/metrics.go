package metrics

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Prometheus metrics for ALL HTTP requests
	HTTPRequestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of handled HTTP requests.",
	})

	// Prometheus metrics for ALL HTTP errors
	HTTPErrorsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "The total number of handled HTTP errors.",
	})

	// Duration of HTTP requests in seconds
	HTTPDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "The duration in seconds of HTTP requests.",
	})
)

func RegisterMetrics() {
	// Register the Prometheus metrics with the global prometheus registry
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPErrorsTotal)
	prometheus.MustRegister(HTTPDuration)
}

func ServeProm(port, path string) error {
	var err error
	// Define Metrics server
	http.Handle(path, promhttp.Handler())

	sm := &http.Server{
		Addr:        port,
		Handler:     promhttp.Handler(),
		ReadTimeout: 10 * time.Second,
	}

	// Start the metrics server
	go func() {
		log.Printf("Starting metrics server on %s", port)
		if err = sm.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()
	return err
}
