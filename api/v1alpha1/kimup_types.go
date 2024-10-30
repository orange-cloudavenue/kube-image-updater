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

type (
	// KimupSpec defines the desired state of Kimup
	KimupSpec struct {
		// +kubebuilder:validation:Required
		// +kubebuilder:description: Kimup instance name
		// The name of the Kimup instance in the suffix of the resource names.
		Name string `json:"name"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Kimup container image
		// Image of the Kimup container. If not set, the default image will be used.
		Image string `json:"image,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Annotations to add to the Kimup pods.
		// Annotations is a key value map that will be added to the Kimup pods.
		Annotations map[string]string `json:"annotations,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Labels to add to the Kimup pods.
		// Labels is a key value map that will be added to the Kimup pods.
		Labels map[string]string `json:"labels,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Env variables to add to the Kimup pods.
		// Env is a list of key value pairs that will be added to the Kimup pods.
		Env []corev1.EnvVar `json:"env,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Kimup container resource limits.
		// Resources is a map of resource requirements that will be added to the Kimup pods.
		Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Affinity settings for the Kimup pods.
		// Affinity is a map of affinity settings that will be added to the Kimup pods.
		Affinity *corev1.Affinity `json:"affinity,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Node selector for the Kimup pods.
		// NodeSelector is a map of node selector settings that will be added to the Kimup pods.
		NodeSelector map[string]string `json:"nodeSelector,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Tolerations for the Kimup pods.
		// Tolerations is a list of tolerations that will be added to the Kimup pods.
		Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: TopologySpreadConstraints for the Kimup pods.
		// TopologySpreadConstraints is a list of constraints that will be added to the Kimup pods.
		TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Service account name for the Kimup pods.
		// +kubebuilder:default:=kimup
		// ServiceAccountName is the name of the service account that will be used by the Kimup pods.
		ServiceAccountName string `json:"serviceAccountName,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Priority class name for the Kimup pods.
		// PriorityClassName is the name of the priority class that will be used by the Kimup pods.
		PriorityClassName string `json:"priorityClassName,omitempty"`

		KimupExtraSpec `json:",inline"`
	}

	// ! Extra

	KimupExtraSpec struct {
		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the metrics settings
		// +kubebuilder:default:={enabled:true}
		// Metrics is a map of settings that will be used to configure the metrics probe. If not set, the probe will be enabled.
		Metrics KimupProbeSpec `json:"metrics,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the healthz settings
		// +kubebuilder:default:={enabled:true}
		// Healthz is a map of settings that will be used to configure the healthz probe. If not set, the probe will be enabled.
		Healthz KimupProbeSpec `json:"healthz,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the log level settings
		// +kubebuilder:default:=info
		// +kubebuilder:validation:Enum=debug;info;warn;error;fatal;panic;trace
		// LogLevel is a string that will be used to configure the log level of the Kimup instance. If not set, the info log level will be used.
		LogLevel string `json:"logLevel,omitempty"`
	}

	KimupProbeSpec struct {
		// +kubebuilder:validation:Optional
		// +kubebuilder:default:=true
		// +kubebuilder:description: Enable or disable the probe
		// Enabled is a boolean that enables or disables the probe. If not set, the probe will be enabled.
		Enabled bool `json:"enabled,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the port settings
		// Port is the port number where the probe will be exposed. If not set, the default port will be used. See https://pkg.go.dev/github.com/orange-cloudavenue/kube-image-updater@v0.0.1/internal/models#pkg-variables.
		Port int32 `json:"port,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Manage the path settings
		// Path is the path where the probe will be exposed. If not set, the default path will be used. See https://pkg.go.dev/github.com/orange-cloudavenue/kube-image-updater@v0.0.1/internal/models#pkg-variables.
		Path string `json:"path,omitempty"`
	}

	// ! Instance

	KimupInstanceSpec struct {
		// +kubebuilder:validation:Required
		// +kubebuilder:description: Kimup instance name
		// The name of the Kimup instance in the suffix of the resource names.
		Name string `json:"name"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Kimup container image
		// Image of the Kimup container. If not set, the default image will be used.
		Image string `json:"image,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Annotations to add to the Kimup pods.
		// Annotations is a key value map that will be added to the Kimup pods.
		Annotations map[string]string `json:"annotations,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Labels to add to the Kimup pods.
		// Labels is a key value map that will be added to the Kimup pods.
		Labels map[string]string `json:"labels,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Env variables to add to the Kimup pods.
		// Env is a list of key value pairs that will be added to the Kimup pods.
		Env []corev1.EnvVar `json:"env,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Kimup container resource limits.
		// Resources is a map of resource requirements that will be added to the Kimup pods.
		Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Affinity settings for the Kimup pods.
		// Affinity is a map of affinity settings that will be added to the Kimup pods.
		Affinity *corev1.Affinity `json:"affinity,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Node selector for the Kimup pods.
		// NodeSelector is a map of node selector settings that will be added to the Kimup pods.
		NodeSelector map[string]string `json:"nodeSelector,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Tolerations for the Kimup pods.
		// Tolerations is a list of tolerations that will be added to the Kimup pods.
		Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: TopologySpreadConstraints for the Kimup pods.
		// TopologySpreadConstraints is a list of constraints that will be added to the Kimup pods.
		TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Service account name for the Kimup pods.
		// +kubebuilder:default:=kimup
		// ServiceAccountName is the name of the service account that will be used by the Kimup pods.
		ServiceAccountName string `json:"serviceAccountName,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description: Priority class name for the Kimup pods.
		// PriorityClassName is the name of the priority class that will be used by the Kimup pods.
		PriorityClassName string `json:"priorityClassName,omitempty"`
	}

	// Status defines the observed state of Kimup
	KimupStatus struct {
		// Status of the Kimup Instance
		// It can be one of the following:
		// - "ready": The kimup instance is ready to serve requests
		// - "resources-created": The Kimup instance resources were created but not yet configured
		State string `json:"state,omitempty"`
	}
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// Kimup is the Schema for the kimups API. Permit to manage the Kimup instances (Controller).
type Kimup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of Kimup
	Spec KimupSpec `json:"spec,omitempty"`

	// Status defines the observed state of Kimup
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
