package models

import "fmt"

var (
	MutatorDefaultPort int = 9443
	MutatorDefaultAddr     = fmt.Sprintf(":%d", MutatorDefaultPort)

	MutatorWebhookConfigurationName = "kimup-mutator"
	MutatorWebhookName              = "image-tag.kimup.cloudavenue.io"
	MutatorServiceName              = MutatorWebhookConfigurationName

	MutatorWebhookPathMutateImageTag = "/mutate/image-tag"
)
