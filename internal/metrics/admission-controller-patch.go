package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	admissionControllerPatch struct{}
)

var (
	// Prometheus metrics
	admissionControllerPatchTotal    prometheus.Counter   = NewCounter("admissionControllerPatch_total", "The total number of patch by the Admission Controller is generate.")
	admissionControllerPatchErrTotal prometheus.Counter   = NewCounter("admissionControllerPatch_error_total", "The total number of patch by the AdmissionController generate with error.")
	admissionControllerPatchDuration prometheus.Histogram = NewHistogram("admissionControllerPatch_duration_seconds", "The duration in seconds of the generated patch by the Admission Controller.")
)

// admissionControllerPatch returns a new admissionControllerPatch.
// This is the metrics for the admissionControllerPatch.
func AdmissionControllerPatch() *admissionControllerPatch {
	return &admissionControllerPatch{}
}

// Total returns the total number of admissionControllerPatch performed.
// The counter is used to observe the number of admissionControllerPatch that have been
// executed. The counter is incremented each time a tag is executed
// A good practice is to use the following pattern:
//
// metrics.admissionControllerPatch().Total().Inc()
func (a *admissionControllerPatch) Total() prometheus.Counter {
	return admissionControllerPatchTotal
}

// TotalErr returns the total number of admissionControllerPatch performed with error.
// The counter is used to observe the number of admissionControllerPatch that failed.
// The counter is incremented each time a tag fails.
// A good practice is to use the following pattern:
//
// metrics.admissionControllerPatch().TotalErr().Inc()
func (a *admissionControllerPatch) TotalErr() prometheus.Counter {
	return admissionControllerPatchErrTotal
}

// Duration returns the duration of the admissionControllerPatch execution.
// A good practice is to use the following pattern:
//
// timeradmissionControllerPatch := metrics.admissionControllerPatch().Duration()
//
// defer timeradmissionControllerPatch.ObserveDuration()
func (a *admissionControllerPatch) Duration() *prometheus.Timer {
	return prometheus.NewTimer(admissionControllerPatchDuration)
}
