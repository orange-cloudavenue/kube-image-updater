package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	// AlertDiscordSpec defines the desired state of AlertDiscord
	AlertDiscordSpec struct {
		// +kubebuilder:validation:Required
		// +kubebuilder:validation:MinItems=1
		// +kubebuilder:validation:example:="123456789012345678"
		ChannelIDs []string `json:"channelIDs"`

		// +kubebuilder:validation:Optional
		CredentialBotToken *ValueOrValueFrom `json:"credentialBotToken,omitempty"`

		// +kubebuilder:validation:Optional
		CredentialOAuth2Token *ValueOrValueFrom `json:"credentialOAuth2,omitempty"`
	}

	// AlertDiscordStatus defines the observed state of AlertDiscord
	AlertDiscordStatus struct {
		AlertCoreStatus `json:",inline"`
	}
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=alertdiscords,scope=Cluster

type AlertDiscord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertDiscordSpec   `json:"spec,omitempty"`
	Status AlertDiscordStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type AlertDiscordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertDiscord `json:"items"`
}

func init() {
	// SchemeBuilder.Register(&AlertDiscord{}, &AlertDiscordList{})
}
