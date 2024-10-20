package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	registry struct {
		RequestTotal      *prometheus.CounterVec `labels:"registry_name" help:"The total number of registry evaluated."`
		RequestErrorTotal *prometheus.CounterVec `labels:"registry_name" help:"The total number of registry evaluated with error."`
		RequestDuration   HistogramVec           `labels:"registry_name" help:"The duration in seconds of registry evaluated."`
	}
)

var registryMetrics registry

// Registry returns a new registry.
// This is the metrics for the registry.
func Registry() registry {
	if registryMetrics.RequestTotal == nil {
		registryMetrics = initMetrics(registry{})
	}

	return registryMetrics
}
