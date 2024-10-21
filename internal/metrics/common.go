package metrics

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type (
	Histogram struct {
		prometheus.Histogram
	}

	HistogramVec struct {
		*prometheus.HistogramVec
	}
)

func (h HistogramVec) NewTimer(labelsValues ...string) *prometheus.Timer {
	return prometheus.NewTimer(h.WithLabelValues(labelsValues...))
}

func (h Histogram) NewTimer() *prometheus.Timer {
	return prometheus.NewTimer(h.Histogram)
}

func initMetrics[T any](s T) T {
	// Reflect the struct type
	t := reflect.TypeOf(s)

	// Get the name of the struct
	name := t.Name()

	// Reflect the struct value
	v := reflect.ValueOf(&s).Elem()

	// Iterate over the struct fields
	for i := 0; i < t.NumField(); i++ {
		// Get the field
		field := t.Field(i)

		// Get the field value
		fieldValue := v.Field(i)

		// Get the help tag
		help := field.Tag.Get("help")

		labels := []string{}
		if field.Tag.Get("labels") != "" {
			// split string by comma
			labels = strings.Split(field.Tag.Get("labels"), ",")
		}

		// Create a new metric
		switch field.Type.String() {
		case "prometheus.Counter", "*prometheus.CounterVec":
			if len(labels) > 0 {
				fieldValue.Set(reflect.ValueOf(newCounterWithVec(buildMetricName(name, field.Name), help, labels)))
			} else {
				fieldValue.Set(reflect.ValueOf(newCounter(buildMetricName(name, field.Name), help)))
			}
		case "metrics.Histogram", "metrics.HistogramVec":
			if len(labels) > 0 {
				fieldValue.Set(reflect.ValueOf(newHistogramVec(buildMetricName(name, field.Name), help, labels)))
			} else {
				fieldValue.Set(reflect.ValueOf(newHistogram(buildMetricName(name, field.Name), help)))
			}
		case "prometheus.Gauge", "*prometheus.GaugeVec":
			if len(labels) > 0 {
				fieldValue.Set(reflect.ValueOf(newGaugeWithVec(buildMetricName(name, field.Name), help, labels)))
			} else {
				fieldValue.Set(reflect.ValueOf(newGauge(buildMetricName(name, field.Name), help)))
			}
		case "prometheus.Summary", "*prometheus.SummaryVec":
			if len(labels) > 0 {
				fieldValue.Set(reflect.ValueOf(newSummaryWithVec(buildMetricName(name, field.Name), help, labels)))
			} else {
				fieldValue.Set(reflect.ValueOf(newSummary(buildMetricName(name, field.Name), help)))
			}
		default:
			panic(fmt.Sprintf("unsupported type %s", field.Type.String()))
		}
	}

	return s
}

func buildMetricName(category, name string) string {
	return fmt.Sprintf(
		"kimup_%s_%s",
		strcase.ToSnake(category),
		strcase.ToSnake(name),
	)
}

// newGauge creates a new Prometheus gauge
// The newGauge use a function to directly register the gauge
// The function returns a prometheus.Gauge
//
// Name: The name of the gauge
// Help: The description help text of the gauge
func newGauge(name, help string) prometheus.Gauge {
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

// newGaugeWithVec creates a new Prometheus gauge with labels
// The newGaugeWithVec use a function to directly register the gauge with labels
// The function returns a prometheus.GaugeVec
//
// Name: The name of the gauge
// Help: The description help text of the gauge
// Labels: The labels of the gauge
func newGaugeWithVec(name, help string, labels []string) *prometheus.GaugeVec {
	if Metrics[MetricTypeGauge] == nil {
		Metrics[MetricTypeGauge] = make(map[string]interface{})
	}

	// Add the gauge to the map
	Metrics[MetricTypeGauge][name] = MetricGaugeVec{
		// Create the metricBase
		metricBase: metricBase{
			Name: name,
			Help: help,
		},
		// Create the gauge prometheus
		GaugeVec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels),
	}

	return Metrics[MetricTypeGauge][name].(MetricGaugeVec).GaugeVec
}

// * Counter

// newCounter creates a new Prometheus counter
// The newCounter use a function to directly register the counter
// The function returns a prometheus.Counter
//
// Name: The name of the counter
// Help: The description help text of the counter
func newCounter(name, help string) prometheus.Counter {
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

// newCounterWithVec creates a new Prometheus counter with labels
// The newCounterWithVec use a function to directly register the counter with labels
// The function returns a prometheus.CounterVec
//
// Name: The name of the counter
// Help: The description help text of the counter
// Labels: The labels of the counter
func newCounterWithVec(name, help string, labels []string) *prometheus.CounterVec {
	if Metrics[MetricTypeCounter] == nil {
		Metrics[MetricTypeCounter] = make(map[string]interface{})
	}

	// Add the counter to the map
	Metrics[MetricTypeCounter][name] = MetricCounterVec{
		// Create the metricBase
		metricBase: metricBase{
			Name: name,
			Help: help,
		},
		// Create the counter prometheus
		CounterVec: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, labels),
	}

	return Metrics[MetricTypeCounter][name].(MetricCounterVec).CounterVec
}

// * Summary

// newSummary creates a new Prometheus summary
// The newSummary use a function to directly register the summary
// The function returns a prometheus.Summary
//
// Name: The name of the summary
// Help: The description help text of the summary
func newSummary(name, help string) prometheus.Summary {
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

// newSummaryWithVec creates a new Prometheus summary with labels
// The newSummaryWithVec use a function to directly register the summary with labels
// The function returns a prometheus.Summary
//
// Name: The name of the summary
// Help: The description help text of the summary
// Labels: The labels of the summary
func newSummaryWithVec(name, help string, labels []string) *prometheus.SummaryVec {
	if Metrics[MetricTypeSummary] == nil {
		Metrics[MetricTypeSummary] = make(map[string]interface{})
	}

	// Add the summary to the map
	Metrics[MetricTypeSummary][name] = MetricSummaryVec{
		// Create the metricBase
		metricBase: metricBase{
			Name: name,
			Help: help,
		},
		// Create the summary prometheus
		SummaryVec: prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: name,
			Help: help,
		}, labels),
	}

	return Metrics[MetricTypeSummary][name].(MetricSummaryVec).SummaryVec
}

// * Histogram

// newHistogram creates a new Prometheus histogram
// The newHistogram use a function to directly register the histogram
// The function returns a prometheus.Histogram
//
// Name: The name of the histogram
// Help: The description help text of the histogram
func newHistogram(name, help string) Histogram {
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
		Histogram: Histogram{
			promauto.NewHistogram(prometheus.HistogramOpts{
				Name: name,
				Help: help,
				// Bucket configuration for microsecond durations
				Buckets: []float64{0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.5, 1, 2, 5, 10},
			}),
		},
	}

	return Metrics[MetricTypeHistogram][name].(MetricHistogram).Histogram
}

// newHistogramVec creates a new Prometheus histogram with labels
// The newHistogramVec use a function to directly register the histogram with labels
// The function returns a prometheus.HistogramVec
//
// Name: The name of the histogram
// Help: The description help text of the histogram
// Labels: The labels of the histogram
func newHistogramVec(name, help string, labels []string) HistogramVec {
	if Metrics[MetricTypeHistogram] == nil {
		Metrics[MetricTypeHistogram] = make(map[string]interface{})
	}

	// Add the histogram to the map
	Metrics[MetricTypeHistogram][name] = MetricHistogramVec{
		// Create the metricBase
		HistogramVec: HistogramVec{
			prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name: name,
				Help: help,
				// Bucket configuration for microsecond durations
				Buckets: []float64{0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.5, 1, 2, 5, 10},
			}, labels),
		},
	}

	return Metrics[MetricTypeHistogram][name].(MetricHistogramVec).HistogramVec
}
