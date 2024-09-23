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

//nolint:gosec
import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
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
		Rules []ImageRule `json:"rules"`
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

	ImageRuleType   string
	ImageActionType string

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
	}
)

const (
	// * ImageRuleType

	// Semver
	ImageRuleTypeSemverMajor ImageRuleType = "semver-major"
	ImageRuleTypeSemverMinor ImageRuleType = "semver-minor"
	ImageRuleTypeSemverPatch ImageRuleType = "semver-patch"

	// Regex
	ImageRuleTypeRegex ImageRuleType = "regex"

	// * ImageActionType

	ImageActionTypeApply           ImageActionType = "apply"
	ImageActionTypeRequestApproval ImageActionType = "request-approval"
	ImageActionTypeNotify          ImageActionType = "notify"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Image is the Schema for the images API
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

// Annotations

// addAnnotation adds an annotation to the Image
func (i *Image) addAnnotation(key, value string) {
	if i.Annotations == nil {
		i.Annotations = map[string]string{}
	}
	i.Annotations[key] = value
}

// ListAnnotations
func (i *Image) ListAnnotations() map[string]string {
	return i.GetAnnotations()
}

// GetAnnotation
func (i *Image) GetAnnotationAction() (annotations.AActionKey, error) {
	action, ok := i.Annotations[annotations.AnnotationActionKey]
	if !ok {
		return "", fmt.Errorf("annotation %s not found", annotations.AnnotationActionKey)
	}

	return annotations.AActionKey(action), nil
}

// SetAnnotationAction
func (i *Image) SetAnnotationAction(action annotations.AActionKey) {
	i.addAnnotation(annotations.AnnotationActionKey, string(action))
}

// getChecksum
func (i *Image) getCheckSum() (string, error) {
	x, err := json.Marshal(i.Spec)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", (md5.Sum(x))), nil //nolint:gosec
}

// GetCheckSum
func (i *Image) GetAnnotationCheckSum() (string, error) {
	checksum, ok := i.Annotations[annotations.AnnotationCheckSumKey]
	if !ok {
		return "", fmt.Errorf("annotation %s not found", annotations.AnnotationCheckSumKey)
	}

	return checksum, nil
}

// RefreshCheckSum
func (i *Image) RefreshCheckSum() error {
	sum, err := i.getCheckSum()
	if err != nil {
		return err
	}

	i.addAnnotation(annotations.AnnotationCheckSumKey, sum)

	return nil
}

// IsChanged
func (i *Image) IsChanged() (bool, error) {
	checksum, err := i.GetAnnotationCheckSum()
	if err != nil {
		return true, err
	}

	sum, err := i.getCheckSum()
	if err != nil {
		return true, err
	}

	return checksum != sum, nil
}
