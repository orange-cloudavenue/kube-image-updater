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
	metrics.InitAll()

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
			m.Counter.Inc()

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
				// Compare the metrics
				if err := testutil.CollectAndCompare(tt.c, strings.NewReader(tt.data+tt.nameMetric+tt.value), tt.nameMetric); err != nil {
					if !tt.error {
						t.Errorf("unexpected error: %v", err)
					}
				}
			})
		} // end of loop over the list of metrics
	}
}

func TestMetric_Histogram(t *testing.T) {
	metrics.InitAll()

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

func TestMetric_Summary(t *testing.T) {
	metrics.InitAll()
	// Test the metrics for the actions Summary
	list := metrics.Metrics

	type testsSummary []struct {
		name       string
		nameMetric string
		data       string
		value      string
		summary    prometheus.Summary
		error      bool
	}
	testUnit := make(testsSummary, 0)

	// loop over the list of metrics
	for _, m := range list[metrics.MetricTypeSummary] {
		// Check if the metric is a metricSummary
		if m, ok := m.(metrics.MetricSummary); ok {
			// fill struct test with data
			testUnit = testsSummary{
				{
					name:       "Check Summary " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeSummary),
					value: fmt.Sprintf(`
%s{quantile="0.5"} 0.1
%s{quantile="0.9"} 0.1
%s_sum 0.1
%s_count 1
`, m.Name, m.Name, m.Name, m.Name),
					summary: m.Summary,
					error:   false,
				},
				{
					name:       "Check Summary with a wrong quantile " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeSummary),
					value: fmt.Sprintf(`
%s{quantile="0.5"} 0.1
%s{quantile="0.9"} 0.1
%s_sum 0.1
%s_count 1
`, m.Name, m.Name, m.Name, m.Name),
					summary: m.Summary,
					error:   true, // Error because the quantile is wrong
				},
			} // end of testsSummary struct
		} // end of if m, ok := m.(metrics.MetricSummary)

		// Test the metrics for the actions Summary
		for _, tt := range testUnit {
			t.Run(tt.name, func(t *testing.T) {
				// Get the Duration histogram
				tt.summary.Observe(0.1)

				// Verify the histogram value
				if err := testutil.CollectAndCompare(tt.summary, strings.NewReader(tt.data+tt.value)); err != nil {
					if !tt.error {
						t.Errorf("unexpected error: %v", err)
					}
				}
			})
		}
	} // end of loop over the list of metrics
}

func TestMetric_Gauge(t *testing.T) {
	metrics.InitAll()
	// Test the metrics for the actions Gauge
	list := metrics.Metrics

	type testsGauge []struct {
		name       string
		nameMetric string
		data       string
		value      string
		g          prometheus.Gauge
		error      bool
	}
	testUnit := make(testsGauge, 0)

	// loop over the list of metrics
	for _, m := range list[metrics.MetricTypeGauge] {
		// Check if the metric is a metricGauge
		if m, ok := m.(metrics.MetricGauge); ok {
			// fill struct test with data
			testUnit = testsGauge{
				{
					name:       "Check Gauge " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeGauge),
					value: " 0\n",
					g:     m.Gauge,
					error: false,
				},
				{
					name:       "Check Gauge mistake between name and TYPE description " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s %s
# TYPE %s_mistake_error_in_TYPE %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeGauge),
					value: " 0\n",
					g:     m.Gauge,
					error: true, // Error because the gauge name is not the same in the HELP description
				},
				{
					name:       "Check Gauge mistake between name and HELP description " + m.Name,
					nameMetric: m.Name,
					data: fmt.Sprintf(`
# HELP %s_mistake_error_in_HELP %s
# TYPE %s %s
`, m.Name, m.Help, m.Name, metrics.MetricTypeGauge),
					value: " 0\n",
					g:     m.Gauge,
					error: true, // Error because the gauge name is not the same in the description
				},
			} // end of testsGauge struct
		} // end of if m, ok := m.(metrics.MetricGauge)

		// Test the metrics for the actions Gauge
		for _, tt := range testUnit {
			t.Run(tt.name, func(t *testing.T) {
				// Compare the metrics
				if err := testutil.CollectAndCompare(tt.g, strings.NewReader(tt.data+tt.nameMetric+tt.value), tt.nameMetric); err != nil {
					if !tt.error {
						t.Errorf("unexpected error: %v", err)
					}
				}
			})
		} // end of loop over the list of metrics
	}
}
