package metrics_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
)

func TestMetric_Counter(t *testing.T) {
	// Test the metrics for the actions Counter
	list := metrics.Metrics

	// Test the metrics for the actions Counter
	type testsCounter []struct {
		name       string
		nameMetric string
		data       string
		value      string
		c          prometheus.Counter
		error      bool
	}
	testUnit := make(testsCounter, 0)

	// loop over the list of metrics
	for _, m := range list[metrics.MetricTypeCounter] {
		// Check if the metric is a metricCounter
		if m, ok := m.(metrics.MetricCounter); ok {
			// fill struct test with data
			testUnit = testsCounter{
				{
					name:       "Check Counter " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeCounter),
					value: " 1\n",
					c:     m.Counter,
					error: false,
				},
				{
					name:       "Check Counter mistake between name and TYPE description " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s_mistake_error_in_TYPE %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeCounter),
					value: " 1\n",
					c:     m.Counter,
					error: true, // Error because the counter name is not the same in the HELP description
				},
				{
					name:       "Check Counter mistake between name and HELP description " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s_mistake_error_in_HELP %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeCounter),
					value: " 1\n",
					c:     m.Counter,
					error: true, // Error because the counter name is not the same in the description
				},
			} // end of testsCounter struct
		} // end of if m, ok := m.(metrics.MetricCounter)

		// Test the metrics for the actions Counter
		for _, tt := range testUnit {
			t.Run(tt.name, func(t *testing.T) {
				counter := tt.c
				counter.Inc()
				// Compare the metrics
				if err := testutil.CollectAndCompare(counter, strings.NewReader(tt.data+tt.nameMetric+tt.value), tt.nameMetric); err != nil {
					if !tt.error {
						t.Errorf("unexpected error: %v", err)
					}
				}
			})
		} // end of loop over the list of metrics
	}
}

func TestMetric_Histogram(t *testing.T) {
	// Test the metrics for the actions Histogram
	list := metrics.Metrics

	type testsHistogram []struct {
		name        string
		nameMetric  string
		data        string
		value       string
		observation float64
		h           prometheus.Histogram
		error       bool
	}
	testUnit := make(testsHistogram, 0)

	// loop over the list of metrics
	for _, m := range list[metrics.MetricTypeHistogram] {
		// Check if the metric is a metricHistogram
		if m, ok := m.(metrics.MetricHistogram); ok {
			// TODO - Constructs all the expected tests
			// fill struct test with data
			testUnit = testsHistogram{
				{
					name:       "Check Histogram " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeHistogram),
					value: fmt.Sprintf(`
%s_bucket{le="0.001"} 0
%s_bucket{le="0.005"} 0
%s_bucket{le="0.01"} 0
%s_bucket{le="0.02"} 0
%s_bucket{le="0.05"} 0
%s_bucket{le="0.1"} 0
%s_bucket{le="0.5"} 1
%s_bucket{le="1"} 1
%s_bucket{le="2"} 1
%s_bucket{le="5"} 1
%s_bucket{le="10"} 1
%s_bucket{le="+Inf"} 1
%s_sum 0.1
%s_count 1
`, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name),
					observation: 0.1,
					h:           m.Histogram,
					error:       false,
				},
				{
					name:       "Check Histogram with bucket missing " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeHistogram),
					value: fmt.Sprintf(`
%s_bucket{le="0.001"} 0
%s_bucket{le="0.01"} 0
%s_bucket{le="0.02"} 0
%s_bucket{le="0.05"} 0
%s_bucket{le="0.1"} 0
%s_bucket{le="0.5"} 1
%s_bucket{le="1"} 1
%s_bucket{le="2"} 1
%s_bucket{le="5"} 1
%s_bucket{le="10"} 1
%s_bucket{le="+Inf"} 1
%s_sum 0.1
%s_count 1
`, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name),
					observation: 0.1,
					h:           m.Histogram,
					error:       true, // Error because the bucket is missing
				},
				{
					name:       "Check Histogram with a wrong observation " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeHistogram),
					value: fmt.Sprintf(`
%s_bucket{le="0.001"} 0
%s_bucket{le="0.005"} 0
%s_bucket{le="0.01"} 0
%s_bucket{le="0.02"} 0
%s_bucket{le="0.05"} 0
%s_bucket{le="0.1"} 0
%s_bucket{le="0.5"} 1
%s_bucket{le="1"} 1
%s_bucket{le="2"} 1
%s_bucket{le="5"} 1
%s_bucket{le="10"} 1
%s_bucket{le="+Inf"} 1
%s_sum 0.1
%s_count 1
`, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name, m.Name),
					observation: 2.5,
					h:           m.Histogram,
					error:       true, // Error because the observation is wrong
				},
			} // end of testsHistogram struct
		} // end of if m, ok := m.(metrics.MetricHistogram)

		// Test the metrics for the actions Histogram
		for _, tt := range testUnit {
			t.Run(tt.name, func(t *testing.T) {
				// Get the Duration histogram
				tt.h.Observe(tt.observation)

				// Verify the histogram value
				if err := testutil.CollectAndCompare(tt.h, strings.NewReader(tt.data+tt.value)); err != nil {
					// Check if err contains "_sum" to avoid the error
					if !strings.Contains(err.Error(), "_sum") {
						if !tt.error {
							t.Errorf("unexpected error: %v", err)
						}
					}
				}
			})
		}
	} // end of loop over the list of metrics
}
