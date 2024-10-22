package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type (
	admissionController struct {
		RequestTotal      prometheus.Counter `help:"The total number of request received."`
		RequestErrorTotal prometheus.Counter `help:"The total number of request received with error."`
		RequestDuration   Histogram          `help:"The duration in seconds of request in admission controller."`
		PatchTotal        prometheus.Counter `help:"The total number of patch action performed."`
		PatchErrorTotal   prometheus.Counter `help:"The total number of patch action performed with error."`
		PatchDuration     Histogram          `help:"The duration in seconds of patch in admission controller."`
	}
)

var admissionControllerMetrics admissionController

// admissionController returns a new admissionController.
// This is the metrics for the admissionController.
func AdmissionController() admissionController {
	if admissionControllerMetrics.RequestTotal == nil {
		admissionControllerMetrics = initMetrics(admissionController{})
	}

	return admissionControllerMetrics
}
