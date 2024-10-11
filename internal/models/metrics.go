package models

import "fmt"

var (
	// Used to enable metrics
	MetricsFlagName = "metrics"

	MetricsPortFlagName       = MetricsFlagName + "-port"
	MetricsDefaultPort  int32 = 9080

	MetricsDefaultAddr = ":" + fmt.Sprintf("%d", MetricsDefaultPort)

	MetricsPathFlagName = MetricsFlagName + "-path"
	MetricsDefaultPath  = "/metrics"
)
