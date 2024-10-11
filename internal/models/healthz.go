package models

import "fmt"

var (
	// Used to enable healthz
	HealthzFlagName = "healthz"

	HealthzPortFlagName       = HealthzFlagName + "-port"
	HealthzDefaultPort  int32 = 9081

	HealthzDefaultAddr = ":" + fmt.Sprintf("%d", HealthzDefaultPort)

	HealthzPathFlagName = HealthzFlagName + "-path"
	HealthzDefaultPath  = "/healthz"
)
