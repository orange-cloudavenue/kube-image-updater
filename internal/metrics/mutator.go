package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	mutator struct {
		RequestTotal      prometheus.Counter `help:"The total number of request received."`
		RequestErrorTotal prometheus.Counter `help:"The total number of request received with error."`
		RequestDuration   Histogram          `help:"The duration in seconds of request in admission controller."`
		PatchTotal        prometheus.Counter `help:"The total number of patch action performed."`
		PatchErrorTotal   prometheus.Counter `help:"The total number of patch action performed with error."`
		PatchDuration     Histogram          `help:"The duration in seconds of patch in admission controller."`
	}
)

var mutatorMetrics mutator

// mutator returns a new mutator.
// This is the metrics for the mutator.
func Mutator() mutator {
	if mutatorMetrics.RequestTotal == nil {
		mutatorMetrics = initMetrics(mutator{})
	}

	return mutatorMetrics
}
