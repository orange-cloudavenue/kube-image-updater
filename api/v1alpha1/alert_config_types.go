package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	AlertConfigSpec struct {
		// +kubebuilder:validation:Optional
		Discord *AlertDiscordSpec `json:"discord,omitempty"`
	}

	// AlertDiscordSpec defines the desired state of AlertDiscord
	AlertDiscordSpec struct {
		// +kubebuilder:validation:Required
		WebhookURL ValueOrValueFrom `json:"webhookURL"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:List of users or roles to notify.
		Mentions []string `json:"mentions,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:Timeout specifies a time limit for the request to be made.
		// +kubebuilder:default:10s
		Timeout string `json:"timeout,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:validation:Schemaless
		// +kubebuilder:validation:Type=string
		TemplateBody string `json:"templateBody,omitempty"`
	}

	// AlertDiscordStatus defines the observed state of AlertDiscord
	AlertConfigStatus struct{}
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=alertconfig,scope=Cluster

type AlertConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertConfigSpec   `json:"spec,omitempty"`
	Status AlertConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type AlertConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertConfig{}, &AlertConfigList{})
}
