package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	rules struct{}
)

var (
	// Prometheus metrics
	rulesTotal    prometheus.Counter   = NewCounter("rules_total", "The total number of rules evaluated.")
	rulesErrTotal prometheus.Counter   = NewCounter("rules_error_total", "The total number of rules evaluated with error.")
	rulesDuration prometheus.Histogram = NewHistogram("rules_duration_seconds", "The duration in seconds of rules evaluated.")
)

// Rules returns a new rules.
// This is the metrics for the rules.
func Rules() *rules {
	return &rules{}
}

// Total returns the total number of rule performed.
// The counter is used to observe the number of rules that have been
// executed. The counter is incremented each time a rule is executed
// A good practice is to use the following pattern:
//
// metrics.Rules().Total().Inc()
func (a *rules) Total() prometheus.Counter {
	return rulesTotal
}

// TotalErr returns the total number of rule performed with error.
// The counter is used to observe the number of rules that failed.
// The counter is incremented each time a rule fails.
// A good practice is to use the following pattern:
//
// metrics.Rules().TotalErr().Inc()
func (a *rules) TotalErr() prometheus.Counter {
	return rulesErrTotal
}

// Duration returns the duration of the rule execution.
// A good practice is to use the following pattern:
//
// timerRules := prometheus.NewTimer(metrics.Rules().Duration())

func (a *rules) Duration() *prometheus.Timer {
	return prometheus.NewTimer(rulesDuration)
}
