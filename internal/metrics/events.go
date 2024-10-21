package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	events struct {
		TriggeredTotal     prometheus.Counter `help:"The total number of events triggered."`
		TriggerdErrorTotal prometheus.Counter `help:"The total number of events triggered with error."`
		TriggeredDuration  Histogram          `help:"The duration in seconds of events triggered."`
	}
)

var eventsMetrics events

// Events returns a new events.
// This is the metrics for the events.
func Events() events {
	if eventsMetrics.TriggeredTotal == nil {
		eventsMetrics = initMetrics(events{})
	}

	return eventsMetrics
}
