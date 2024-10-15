/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type (
	// KimupSpec defines the desired state of Kimup
	KimupSpec struct {
		// TODO add namespace and serviceaccount settings

		// +kubebuilder:validation:Optional
		Controller *KimupControllerSpec `json:"controller"`

		// +kubebuilder:validation:Optional
		AdmissionController *KimupAdmissionControllerSpec `json:"admissionController"`
	}

	// ! Controller

	KimupControllerSpec struct {
		KimupInstanceSpec `json:",inline"`

		KimupExtraSpec `json:",inline"`

		// Service *KimupServiceSpec `json:"service,omitempty"`
	}

	// ! AdmissionController

	KimupAdmissionControllerSpec struct {
		// +kubebuilder:validation:Optional
		// +kubebuilder:default:=Deployment
		// +kubebuilder:validation:Enum=Deployment;DaemonSet
		DeploymentType string `json:"deploymentType,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:default:=3
		// +kubebuilder:description: Number of replicas (default: 3) for the admissionController deployment. (Only for Deployment)
		Replicas int32 `json:"replicas,omitempty"`

		KimupInstanceSpec `json:",inline"`

		KimupExtraSpec `json:",inline"`
	}

	// ! Extra

	KimupExtraSpec struct {
		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the metrics settings
		// +kubebuilder:default:={enabled:true}
		Metrics KimupProbeSpec `json:"metrics,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the healthz settings
		// +kubebuilder:default:={enabled:true}
		Healthz KimupProbeSpec `json:"healthz,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the log level settings
		// +kubebuilder:default:=info
		// +kubebuilder:validation:Enum=debug;info;warn;error
		LogLevel string `json:"logLevel,omitempty"`
	}

	KimupProbeSpec struct {
		// +kubebuilder:validation:Optional
		// +kubebuilder:default:=true
		Enabled bool `json:"enabled,omitempty"`

		// +kubebuilder:validation:Optional
		Port int32 `json:"port,omitempty"`

		// +kubebuilder:validation:Optional
		Path string `json:"path,omitempty"`
	}

	// ! Instance

	KimupInstanceSpec struct {
		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Kimup container image
		Image string `json:"image,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Annotations to add to the Kimup pods.
		Annotations map[string]string `json:"annotations,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Labels to add to the Kimup pods.
		Labels map[string]string `json:"labels,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Env variables to add to the Kimup pods.
		Env []corev1.EnvVar `json:"env,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Kimup container resource limits.
		Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Affinity settings for the Kimup pods.
		Affinity *corev1.Affinity `json:"affinity,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Node selector for the Kimup pods.
		NodeSelector map[string]string `json:"nodeSelector,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Tolerations for the Kimup pods.
		Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: TopologySpreadConstraints for the Kimup pods.
		TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Service account name for the Kimup pods.
		// +kubebuilder:default:=kimup
		ServiceAccountName string `json:"serviceAccountName,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Priority class name for the Kimup pods.
		PriorityClassName string `json:"priorityClassName,omitempty"`
	}

	// ! Service

	// KimupServiceSpec struct {
	// 	// +kubebuilder:validation:Optional
	// 	// +kubebuilder:description: Type of the Kimup service
	// 	Type corev1.ServiceType `json:"type,omitempty"`

	// 	// +kubebuilder:validation:Optional
	// 	// +kubebuilder:description: Name of the Kimup service
	// 	Name string `json:"name,omitempty"`

	// 	// +kubebuilder:validation:Optional
	// 	// +kubebuilder:description: Port for the Kimup service
	// 	Port int32 `json:"port,omitempty"`

	// 	// +kubebuilder:validation:Optional
	// 	// +kubebuilder:description: Annotations to add to the Kimup service.
	// 	Annotations map[string]string `json:"annotations,omitempty"`

	// 	// +kubebuilder:validation:Optional
	// 	// +kubebuilder:description: Labels to add to the Kimup service.
	// 	Labels map[string]string `json:"labels,omitempty"`
	// }

	// KimupStatus defines the observed state of Kimup
	KimupStatus struct {
		Controller KimupInstanceStatus `json:"controller,omitempty"`

		AdmissionController KimupInstanceStatus `json:"admissionController,omitempty"`
	}

	KimupInstanceStatus struct {
		// Status of the Kimup Instance
		// It can be one of the following:
		// - "ready": The kimup instance is ready to serve requests
		// - "resources-created": The Kimup instance resources were created but not yet configured
		State string `json:"state,omitempty"`

		// IsRollingUpdate is true if the kimup instance is being updated
		IsRollingUpdate bool `json:"isRollingUpdate,omitempty"`
	}
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Kimup is the Schema for the kimups API
// +kubebuilder:printcolumn:name="Controller",type=string,JSONPath=`.status.controller.state`
// +kubebuilder:printcolumn:name="AdmissionController",type=string,JSONPath=`.status.admissionController.state`
type Kimup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KimupSpec   `json:"spec,omitempty"`
	Status KimupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KimupList contains a list of Kimup
type KimupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kimup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kimup{}, &KimupList{})
}
