package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	events struct{}
)

var (
	// Prometheus metrics
	eventsTotal    prometheus.Counter   = NewCounter("events_total", "The total number of events.")
	eventsTotalErr prometheus.Counter   = NewCounter("events_total_err", "The total number of events with error.")
	eventsDuration prometheus.Histogram = NewHistogram("events_duration_seconds", "The duration in seconds of events.")
)

// Events returns a new events.
// This is the metrics for the events.
func Events() *events {
	return &events{}
}

// Total returns the total number of event performed.
// The counter is used to observe the number of events that have been executed.
// The counter is incremented each time an event is executed
// A good practice is to use the following pattern:
//
// metrics.Events().Total().Inc()
func (a *events) Total() prometheus.Counter {
	return eventsTotal
}

// TotalErr returns the total number of event performed with error.
// The counter is used to observe the number of events that failed.
// The counter is incremented each time an event fails.
// A good practice is to use the following pattern:
//
// metrics.Events().TotalErr().Inc()
func (a *events) TotalErr() prometheus.Counter {
	return eventsTotalErr
}

// Duration returns a prometheus histogram object.
// The histogram is used to observe the duration of the events execution.
// A good practice is to use the following pattern:
//
// timerEvents := metrics.Events().Duration()
//
// defer timerEvents.ObserveDuration()
func (a *events) Duration() *prometheus.Timer {
	return prometheus.NewTimer(eventsDuration)
}
