package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	actions struct {
		ExecutedTotal      prometheus.Counter `help:"The total number of action performed."`
		ExecutedErrorTotal prometheus.Counter `help:"The total number of action performed with error."`
		ExecutedDuration   Histogram          `help:"The duration in seconds of action performed."`
	}
)

var actionsMetrics actions

// Actions returns a new actions.
// This is the metrics for the actions.
func Actions() actions {
	if actionsMetrics.ExecutedTotal == nil {
		actionsMetrics = initMetrics(actions{})
	}

	return actionsMetrics
}
