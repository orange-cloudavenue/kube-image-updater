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
	// ImageSpec defines the desired state of Image
	ImageSpec struct {
		// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
		// Important: Run "make" to regenerate code after modifying this file

		// +kubebuilder:validation:Required
		Image string `json:"image"`

		// +kubebuilder:validation:Optional
		ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:default:="latest"
		// +kubebuilder:example:="v1.2.0"
		BaseTag string `json:"baseTag,omitempty"`

		// +kubebuilder:validation:Required
		// +kubebuilder:validation:MinItems=1
		Triggers []ImageTrigger `json:"triggers"`

		// +kubebuilder:validation:Required
		// +kubebuilder:validation:MinItems=1
		Rules []ImageRule `json:"rules"`
	}

	// ImageTrigger
	ImageTrigger struct {
		// +kubebuilder:validation:Required
		// +kubebuilder:validation:Enum=crontab;webhook
		Type ImageTriggerType `json:"type"`

		// +kubebuilder:validation:Optional
		Value string `json:"value"`
	}

	// ImageRule
	ImageRule struct {
		// +kubebuilder:validation:Required
		Name string `json:"name"`

		// +kubebuilder:validation:Required
		// +kubebuilder:validation:Enum=semver-major;semver-minor;semver-patch;regex
		Type ImageRuleType `json:"type"`

		// +kubebuilder:validation:Optional
		Value string `json:"value,omitempty"`

		// +kubebuilder:validation:Required
		// +kubebuilder:validation:MinItems=1
		Actions []ImageAction `json:"actions"`
	}

	ImageTriggerType string
	ImageRuleType    string
	ImageActionType  string

	// ImageAction
	ImageAction struct {
		// +kubebuilder:validation:Required
		// +kubebuilder:validation:Enum=apply;request-approval;notify
		Type ImageActionType `json:"type"`
	}

	// ImageStatus defines the observed state of Image
	ImageStatus struct {
		// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
		// Important: Run "make" to regenerate code after modifying this file
		Tag string `json:"tag"`
	}
)

const (
	// * ImageRuleType

	// TODO(mickael) use const in rules package
	// Semver
	ImageRuleTypeSemverMajor ImageRuleType = "semver-major"
	ImageRuleTypeSemverMinor ImageRuleType = "semver-minor"
	ImageRuleTypeSemverPatch ImageRuleType = "semver-patch"

	// Regex
	ImageRuleTypeRegex ImageRuleType = "regex"

	// TODO(mickael) use const in action package
	// * ImageActionType
	ImageActionTypeApply           ImageActionType = "apply"
	ImageActionTypeRequestApproval ImageActionType = "request-approval"
	ImageActionTypeNotify          ImageActionType = "notify"

	// TODO(mickael) use const in trigger package
	// * ImageTriggerType
	ImageTriggerTypeCrontab ImageTriggerType = "crontab"
	ImageTriggerTypeWebhook ImageTriggerType = "webhook"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Image is the Schema for the images API
// +kubebuilder:printcolumn:name="Tag",type=string,JSONPath=`.status.tag`
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ImageList contains a list of Image
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Image{}, &ImageList{})
}

// SetStatusTag sets the status tag of the image
func (i *Image) SetStatusTag(tag string) {
	i.Status.Tag = tag
}
