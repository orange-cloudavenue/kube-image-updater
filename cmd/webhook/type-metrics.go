package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Prometheus metrics for ALL HTTP requests
	totalHTTPRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_http_requests",
		Help: "The total number of handled HTTP requests.",
	})

	// Prometheus metrics for ALL HTTP errors
	totalHTTPErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "total_http_errors",
		Help: "The total number of handled HTTP errors.",
	})

	// Duration of HTTP requests in seconds
	httpDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
	})
	// promauto.NewHistogramVec(prometheus.HistogramOpts{
	// 	Name: "http_response_time_seconds",
	// 	Help: "Duration of HTTP requests.",
	// })
)
