package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	rules struct {
		EvaluatedTotal      prometheus.Counter `help:"The total number of rules evaluated."`
		EvaluatedErrorTotal prometheus.Counter `help:"The total number of rules evaluated with error."`
		EvaluatedDuration   Histogram          `help:"The duration in seconds of rules evaluated."`
	}
)

var ruleMetrics rules

// Rules returns a new rules.
// This is the metrics for the rules.
func Rules() rules {
	if ruleMetrics.EvaluatedTotal == nil {
		ruleMetrics = initMetrics(rules{})
	}

	return ruleMetrics
}
