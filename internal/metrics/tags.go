package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	tags struct {
		TagsAvailableSum  *prometheus.SummaryVec `labels:"image_name" help:"The total number of tags available for an image."`
		RequestTotal      prometheus.Counter     `help:"The total number of requests to list tags."`
		RequestErrorTotal prometheus.Counter     `help:"The total number returned an error when calling list tags."`
		RequestDuration   Histogram              `help:"The duration in seconds of the request to list tags."`
	}
)

var tagMetrics tags

// Tags returns a new tags.
// This is the metrics for the tags.
func Tags() tags {
	if tagMetrics.TagsAvailableSum == nil {
		tagMetrics = initMetrics(tags{})
	}

	return tagMetrics
}
