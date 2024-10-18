package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	admissionController struct{}
)

var (
	// Prometheus metrics
	admissionControllerTotal    prometheus.Counter   = NewCounter("admissionController_total", "The total number of action performed.")
	admissionControllerTotalErr prometheus.Counter   = NewCounter("admissionController_total_err", "The total number of action performed with error.")
	admissionControllerDuration prometheus.Histogram = NewHistogram("admissionController_duration_seconds", "The duration in seconds of action performed.")
)

// admissionController returns a new admissionController.
// This is the metrics for the admissionController.
func AdmissionController() *admissionController {
	return &admissionController{}
}

// Total returns the total number of admission controller is performed.
// The counter is used to observe the number of admissionController that have been executed.
// The counter is incremented each time an admission controller is executed
// A good practice is to use the following pattern:
//
// metrics.admissionController().Total().Inc()
func (a *admissionController) Total() prometheus.Counter {
	return admissionControllerTotal
}

// TotalErr returns the total number of admission controller performed with error.
// The counter is used to observe the number of admissionController that failed.
// The counter is incremented each time an admission controller fails.
// A good practice is to use the following pattern:
//
// metrics.admissionController().TotalErr().Inc()
func (a *admissionController) TotalErr() prometheus.Counter {
	return admissionControllerTotalErr
}

// ExecuteDuration returns the duration of the admission controller execution.
// A good practice is to use the following pattern:
//
// timer := metrics.AdmissionController().Duration()
//
// defer timer.ObserveDuration()
func (a *admissionController) Duration() *prometheus.Timer {
	return prometheus.NewTimer(admissionControllerDuration)
}
