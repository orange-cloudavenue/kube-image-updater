package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
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

	MetricCounterVec struct {
		metricBase
		CounterVec *prometheus.CounterVec
	}

	MetricGauge struct {
		metricBase
		Gauge prometheus.Gauge
	}

	MetricGaugeVec struct {
		metricBase
		GaugeVec *prometheus.GaugeVec
	}

	MetricHistogram struct {
		metricBase
		Histogram Histogram
	}

	MetricHistogramVec struct {
		metricBase
		HistogramVec HistogramVec
	}

	MetricSummary struct {
		metricBase
		Summary prometheus.Summary
	}

	MetricSummaryVec struct {
		metricBase
		SummaryVec *prometheus.SummaryVec
	}

	MetricType string

	Metric interface {
		GetHelp() string
		GetName() string
	}
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

func InitAll() {
	Actions()
	Events()
	Tags()
	Rules()
	Registry()
	AdmissionController()
}

// GetHelp returns the help text of the metric
func (m metricBase) GetHelp() string {
	return m.Help
}

// GetName returns the name of the metric
func (m metricBase) GetName() string {
	return m.Name
}
