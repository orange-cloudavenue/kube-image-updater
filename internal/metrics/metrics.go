package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// metricBase is a base struct for all metrics
type (
	metricBase struct {
		Name string
		Help string
	}

	MetricCounter struct {
		metricBase
		Counter prometheus.Counter
	}

	MetricGauge struct {
		metricBase
		Gauge prometheus.Gauge
	}

	MetricHistogram struct {
		metricBase
		Histogram prometheus.Histogram
	}

	MetricSummary struct {
		metricBase
		Summary prometheus.Summary
	}

	MetricType string
)

const (
	// MetricTypeCounter is the type of the metric counter
	MetricTypeCounter MetricType = "counter"
	// MetricTypeGauge is the type of the metric gauge
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeHistogram is the type of the metric histogram
	MetricTypeHistogram MetricType = "histogram"
	// MetricTypeSummary is the type of the metric summary
	MetricTypeSummary MetricType = "summary"
)

var Metrics = make(map[MetricType]map[string]interface{})

// NewCounter creates a new Prometheus counter
// The NewCounter use a function to directly register the counter
// The function returns a prometheus.Counter
//
// Name: The name of the counter
// Help: The description help text of the counter
func NewCounter(name, help string) prometheus.Counter {
	if Metrics[MetricTypeCounter] == nil {
		Metrics[MetricTypeCounter] = make(map[string]interface{})
	}

	// Add the counter to the map
	Metrics[MetricTypeCounter][name] = MetricCounter{
		// Create the metricBase
		metricBase: metricBase{
			Name: name,
			Help: help,
		},
		// Create the counter prometheus
		Counter: promauto.NewCounter(prometheus.CounterOpts{
			Name: name,
			Help: help,
		}),
	}

	return Metrics[MetricTypeCounter][name].(MetricCounter).Counter
}

// NewGauge creates a new Prometheus gauge
// The NewGauge use a function to directly register the gauge
// The function returns a prometheus.Gauge
//
// Name: The name of the gauge
// Help: The description help text of the gauge
func NewGauge(name, help string) prometheus.Gauge {
	if Metrics[MetricTypeGauge] == nil {
		Metrics[MetricTypeGauge] = make(map[string]interface{})
	}

	// Add the gauge to the map
	Metrics[MetricTypeGauge][name] = MetricGauge{
		// Create the metricBase
		metricBase: metricBase{
			Name: name,
			Help: help,
		},
		// Create the gauge prometheus
		Gauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}),
	}

	return Metrics[MetricTypeGauge][name].(MetricGauge).Gauge
}

// NewHistogram creates a new Prometheus histogram
// The NewHistogram use a function to directly register the histogram
// The function returns a prometheus.Histogram
//
// Name: The name of the histogram
// Help: The description help text of the histogram
func NewHistogram(name, help string) prometheus.Histogram {
	if Metrics[MetricTypeHistogram] == nil {
		Metrics[MetricTypeHistogram] = make(map[string]interface{})
	}

	// Add the histogram to the map
	Metrics[MetricTypeHistogram][name] = MetricHistogram{
		// Create the metricBase
		metricBase: metricBase{
			Name: name,
			Help: help,
		},
		// Create the histogram prometheus
		Histogram: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: name,
			Help: help,
			// Bucket configuration for microsecond durations
			Buckets: []float64{0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		}),
	}

	return Metrics[MetricTypeHistogram][name].(MetricHistogram).Histogram
}

// NewSummary creates a new Prometheus summary
// The NewSummary use a function to directly register the summary
// The function returns a prometheus.Summary
//
// Name: The name of the summary
// Help: The description help text of the summary
func NewSummary(name, help string) prometheus.Summary {
	if Metrics[MetricTypeSummary] == nil {
		Metrics[MetricTypeSummary] = make(map[string]interface{})
	}

	// Add the summary to the map
	Metrics[MetricTypeSummary][name] = MetricSummary{
		// Create the metricBase
		metricBase: metricBase{
			Name: name,
			Help: help,
		},
		// Create the summary prometheus
		Summary: promauto.NewSummary(prometheus.SummaryOpts{
			Name: name,
			Help: help,
		}),
	}

	return Metrics[MetricTypeSummary][name].(MetricSummary).Summary
}
