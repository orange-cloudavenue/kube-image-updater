package models

import "fmt"

var (
	AdmissionControllerDefaultPort int32 = 9099
	AdmissionControllerDefaultAddr       = fmt.Sprintf(":%d", AdmissionControllerDefaultPort)

	AdmissionControllerMutatingWebhookConfigurationName = "kimup-admission-controller-mutating"
	AdmissionControllerMutatingWebhookName              = "image-tag.kimup.io"
	AdmissionControllerServiceName                      = AdmissionControllerMutatingWebhookConfigurationName

	AdmissionControllerWebhookPathMutateImageTag = "/mutate/image-tag"
)
