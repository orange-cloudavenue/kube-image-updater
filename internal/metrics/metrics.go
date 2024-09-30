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
	MyHTTPDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "my_http_response_time_seconds",
		Help: "The duration in seconds of HTTP requests.",
	})
)

func RegisterMetrics() (err error) {
	// Register the Prometheus metrics with the global prometheus registry
	if err = prometheus.Register(HTTPRequestsTotal); err != nil {
		return err
	}
	if err = prometheus.Register(HTTPErrorsTotal); err != nil {
		return err
	}
	if err = prometheus.Register(MyHTTPDuration); err != nil {
		log.Fatalf("Failed to register metrics: %v", MyHTTPDuration)
		return err
	}
	return nil
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
