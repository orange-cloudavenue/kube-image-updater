package models

import "fmt"

var (
	MutatorDefaultPort int32 = 8443
	MutatorDefaultAddr       = fmt.Sprintf(":%d", MutatorDefaultPort)

	MutatorMutatingWebhookConfigurationName = "kimup-admission-controller-mutating"
	MutatorMutatingWebhookName              = "image-tag.kimup.io"
	MutatorServiceName                      = MutatorMutatingWebhookConfigurationName

	MutatorWebhookPathMutateImageTag = "/mutate/image-tag"
)
