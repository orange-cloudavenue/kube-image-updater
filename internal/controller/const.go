package controller

var Version = "dev"

const (
	BaseKimupImage = "ghcr.io/orange-cloudavenue/"

	KimupName = "kimup"

	KimupOperator      = "operator"
	KimupOperatorName  = KimupName + "-" + KimupOperator
	KimupOperatorImage = BaseKimupImage + KimupOperatorName

	KimupController      = "controller"
	KimupControllerName  = KimupName + "-" + KimupController
	KimupControllerImage = BaseKimupImage + KimupControllerName

	KimupAdmissionController      = "admission-controler"
	KimupAdmissionControllerName  = KimupName + "-" + KimupAdmissionController
	KimupAdmissionControllerImage = BaseKimupImage + KimupAdmissionControllerName
)

const (
	StateResourcesCreated string = "resources-created"

	StateReady string = "ready"
)

const (
	// Recommended Kubernetes Application Labels
	// KubernetesAppNameLabel is the name of the application
	KubernetesAppNameLabelKey = "app.kubernetes.io/name"

	// KubernetesAppVersionLabel is the version of the application
	KubernetesAppVersionLabelKey = "app.kubernetes.io/version"

	// KubernetesAppComponentLabel is the component of the application
	KubernetesAppComponentLabelKey = "app.kubernetes.io/component"

	KubernetesAppInstanceNameLabel = "app.kubernetes.io/instance"

	// KubernetesManagedByLabel is the tool being used to manage the operation of an application
	KubernetesManagedByLabelKey = "app.kubernetes.io/managed-by"

	// KubernetesPartOfLabel is the name of a higher level application this one is part of
	KubernetesPartOfLabelKey = "app.kubernetes.io/part-of"
)
