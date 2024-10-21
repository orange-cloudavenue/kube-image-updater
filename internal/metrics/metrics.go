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
)

const (
	// MetricTypeCounter is the type of the metric counter
	MetricTypeCounter MetricType = "counter"
	// MetricTypeCounterVec is the type of the metric counter with labels
	MetricTypeCounterVec MetricType = "counterVec"
	// MetricTypeGauge is the type of the metric gauge
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeGaugeVec is the type of the metric gauge with labels
	MetricTypeGaugeVec MetricType = "gaugeVec"
	// MetricTypeHistogram is the type of the metric histogram
	MetricTypeHistogram MetricType = "histogram"
	// MetricTypeHistogramVec is the type of the metric histogram with labels
	MetricTypeHistogramVec MetricType = "histogramVec"
	// MetricTypeSummary is the type of the metric summary
	MetricTypeSummary MetricType = "summary"
	// MetricTypeSummaryVec is the type of the metric summary with labels
	MetricTypeSummaryVec MetricType = "summaryVec"
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
