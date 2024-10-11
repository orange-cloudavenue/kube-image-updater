package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	tags struct{}
)

var (
	// Prometheus metrics
	tagsTotal    prometheus.Counter   = NewCounter("tags_total", "The total number of func tags is called to list tags.")
	tagsErrTotal prometheus.Counter   = NewCounter("tags_error_total", "The total number return by the func tags with error.")
	tagsDuration prometheus.Histogram = NewHistogram("tags_duration_seconds", "The duration in seconds for func tags to list the tags.")
)

// Tags returns a new tags.
// This is the metrics for the tags.
func Tags() *tags {
	return &tags{}
}

// Total returns the total number of func tags is called.
// The counter is used to observe the number of func tags is executed.
// The counter is incremented each time a tag is executed
// A good practice is to use the following pattern:
//
// metrics.Tags().Total().Inc()
func (a *tags) Total() prometheus.Counter {
	return tagsTotal
}

// TotalErr returns the total number of func tags called with error.
// The counter is used to observe the number of func tags that failed.
// The counter is incremented each time a tag fails.
// A good practice is to use the following pattern:
//
// metrics.Tags().TotalErr().Inc()
func (a *tags) TotalErr() prometheus.Counter {
	return tagsErrTotal
}

// Duration returns the duration of the func tags execution.
// A good practice is to use the following pattern:
//
// timerTags := metrics.Tags().Duration()
//
// defer timerTags.ObserveDuration()
func (a *tags) Duration() *prometheus.Timer {
	return prometheus.NewTimer(tagsDuration)
}
