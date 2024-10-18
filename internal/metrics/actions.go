package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	actions struct{}
)

var (
	// Prometheus metrics
	actionsTotal    prometheus.Counter   = NewCounter("actions_total", "The total number of action performed.")
	actionsErrTotal prometheus.Counter   = NewCounter("actions_error_total", "The total number of action performed with error.")
	actionsDuration prometheus.Histogram = NewHistogram("actions_duration_seconds", "The duration in seconds of action performed.")
)

// Actions returns a new actions.
// This is the metrics for the actions.
func Actions() *actions {
	return &actions{}
}

// Total returns the total number of action performed.
// The counter is used to observe the number of actions that have been executed.
// The counter is incremented each time an action is executed
// A good practice is to use the following pattern:
//
// metrics.Actions().Total().Inc()
func (a *actions) Total() prometheus.Counter {
	return actionsTotal
}

// TotalErr returns the total number of action performed with error.
// The counter is used to observe the number of actions that failed.
// The counter is incremented each time an action fails.
// A good practice is to use the following pattern:
//
// metrics.Actions().TotalErr().Inc()
func (a *actions) TotalErr() prometheus.Counter {
	return actionsErrTotal
}

// ExecuteDuration returns the duration of the action execution.
// A good practice is to use the following pattern:
//
// timerActions := metrics.Actions().Duration()
//
// defer timerActions.ObserveDuration()
func (a *actions) Duration() *prometheus.Timer {
	return prometheus.NewTimer(actionsDuration)
}
