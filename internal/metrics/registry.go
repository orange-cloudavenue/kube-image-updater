package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	registry struct{}
)

var (
	// Prometheus metrics
	registryTotal    prometheus.Counter   = NewCounter("registry_total", "The total number of registry evaluated.")
	registryErrTotal prometheus.Counter   = NewCounter("registry_error_total", "The total number of registry evaluated with error.")
	registryDuration prometheus.Histogram = NewHistogram("registry_duration_seconds", "The duration in seconds of registry evaluated.")
)

// Registry returns a new registry.
// This is the metrics for the registry.
func Registry() *registry {
	return &registry{}
}

// Total returns the total number of registry is called.
// The counter is used to observe the number of registry that have been executed.
// The counter is incremented each time an registry is executed
// A good practice is to use the following pattern:
//
// metrics.Registry().Total().Inc()
func (a *registry) Total() prometheus.Counter {
	return registryTotal
}

// TotalErr returns the total number of registry is called with error.
// The counter is used to observe the number of registry that failed.
// The counter is incremented each time an registry fails.
// A good practice is to use the following pattern:
//
// metrics.Registry().TotalErr().Inc()
func (a *registry) TotalErr() prometheus.Counter {
	return registryErrTotal
}

// Duration returns the duration of the registry execution.
// A good practice is to use the following pattern:
//
// timerRegistry := metrics.Registry().Duration()

// defer timerRegistry.ObserveDuration()
func (a *registry) Duration() *prometheus.Timer {
	return prometheus.NewTimer(registryDuration)
}
