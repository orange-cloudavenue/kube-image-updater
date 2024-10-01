package metrics

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	webhookMetricsPath string = "/metrics"
	webhookMetricsPort string = ":9080"
)

// NewCounter creates a new Prometheus counter
// The NewCounter use a function to directly register the counter
// The function returns a prometheus.Counter
//
// Name: The name of the counter
// Help: The description help text of the counter
func NewCounter(name, help string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
}

// NewGauge creates a new Prometheus gauge
// The NewGauge use a function to directly register the gauge
// The function returns a prometheus.Gauge
//
// Name: The name of the gauge
// Help: The description help text of the gauge
func NewGauge(name, help string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
}

// NewHistogram creates a new Prometheus histogram
// The NewHistogram use a function to directly register the histogram
// The function returns a prometheus.Histogram
//
// Name: The name of the histogram
// Help: The description help text of the histogram
func NewHistogram(name, help string) prometheus.Histogram {
	return promauto.NewHistogram(prometheus.HistogramOpts{
		Name: name,
		Help: help,
	})
}

// NewSummary creates a new Prometheus summary
// The NewSummary use a function to directly register the summary
// The function returns a prometheus.Summary
//
// Name: The name of the summary
// Help: The description help text of the summary
func NewSummary(name, help string) prometheus.Summary {
	return promauto.NewSummary(prometheus.SummaryOpts{
		Name: name,
		Help: help,
	})
}

func init() {
	flag.StringVar(&webhookMetricsPort, "metrics-port", webhookMetricsPort, "Metrics server port. ex: :9080")
	flag.StringVar(&webhookMetricsPath, "metrics-path", webhookMetricsPath, "Metrics server path. ex: /metrics")
}

// ServeProm starts a Prometheus metrics server
// TODO - Add context to cancel the server
// in order to stop the server gracefully
func ServeProm() error {
	var err error
	// Define Metrics server
	mux := http.NewServeMux()
	mux.Handle(webhookMetricsPath, promhttp.Handler())

	sm := &http.Server{
		Addr:        webhookMetricsPort,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}

	// Start the metrics server
	go func() {
		log.Printf("Starting metrics server on %s", webhookMetricsPort)
		if err = sm.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()
	return err
}
