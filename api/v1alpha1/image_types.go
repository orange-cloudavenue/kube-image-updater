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
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
	"github.com/orange-cloudavenue/kube-image-updater/internal/triggers"
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
		// +kubebuilder:default:=false
		// +kubebuilder:example:=true
		InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`

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
		Type triggers.Name `json:"type"`

		// +kubebuilder:validation:Optional
		Value string `json:"value"`
	}

	// ImageRule
	ImageRule struct {
		// +kubebuilder:validation:Required
		Name string `json:"name"`

		// +kubebuilder:validation:Required
		// +kubebuilder:validation:Enum=calver-major;calver-minor;calver-patch;calver-prerelease;semver-major;semver-minor;semver-patch;regex;always
		Type rules.Name `json:"type"`

		// +kubebuilder:validation:Optional
		Value string `json:"value,omitempty"`

		// +kubebuilder:validation:Required
		// +kubebuilder:validation:MinItems=1
		Actions []ImageAction `json:"actions"`
	}

	// ImageAction
	ImageAction struct {
		// +kubebuilder:validation:Required
		// +kubebuilder:validation:Enum=apply;request-approval;alert-discord
		Type string `json:"type"`

		// +kubebuilder:validation:Optional
		Data ValueOrValueFrom `json:"data,omitempty"`
	}

	// ImageStatus defines the observed state of Image
	ImageStatus struct {
		// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
		// Important: Run "make" to regenerate code after modifying this file
		Tag    string              `json:"tag"`
		Result ImageStatusLastSync `json:"result"`
		Time   string              `json:"time"`
	}
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Image is the Schema for the images API
// +kubebuilder:printcolumn:name="Image",type=string,JSONPath=`.spec.image`
// +kubebuilder:printcolumn:name="Tag",type=string,JSONPath=`.status.tag`
// +kubebuilder:printcolumn:name="Last-Result",type=string,JSONPath=`.status.result`
// +kubebuilder:printcolumn:name="Last-Sync",type=date,JSONPath=`.status.time`
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

// SetStatusExecution sets the status execution of the image
func (i *Image) SetStatusResult(result ImageStatusLastSync) {
	i.Status.Result = result
}

// SetStatusTime sets the status time of the image
func (i *Image) SetStatusTime(time string) {
	i.Status.Time = time
}

// GetImageWithTag returns the image name with the tag
func (i *Image) GetImageWithTag() string {
	if i.Status.Tag == "" {
		return i.Spec.Image + ":" + i.Spec.BaseTag
	}

	return i.Spec.Image + ":" + i.Status.Tag
}

// GetImageWithoutTag returns the image name without the tag
func (i *Image) GetImageWithoutTag() string {
	return i.Spec.Image
}

// Get Tag returns the tag of the image
func (i *Image) GetTag() string {
	if i.Status.Tag == "" {
		return i.Spec.BaseTag
	}

	return i.Status.Tag
}

// ImageIsEqual checks if the provided image string is equal to the image
// specified in the Image struct. The provided image string can be in the
// format of "image:tag" or just "image". This function compares only the
// image name, ignoring any tags. It returns true if the image names are
// equal, and false otherwise.
func (i *Image) ImageIsEqual(image string) bool {
	// Define if image has a tag
	// image possible format are:
	// - image:tag
	// - image

	// Split image and tag
	imageSplit := strings.Split(image, ":")
	imageName := imageSplit[0]

	return i.Spec.Image == imageName
}
