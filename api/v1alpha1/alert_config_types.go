package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	AlertConfigSpec struct {
		// +kubebuilder:validation:Optional
		Discord *AlertDiscordSpec `json:"discord,omitempty"`

		// +kubebuilder:validation:Optional
		Email *AlertEmailSpec `json:"email,omitempty"`
	}

	// AlertEmailSpec defines the desired state of AlertEmail
	AlertEmailSpec struct {
		// +kubebuilder:validation:Required
		// +kubebuilder:description:Host specifies the SMTP server to connect to.
		Host ValueOrValueFrom `json:"host"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:Port specifies the port to connect to the SMTP server.
		// +kubebuilder:default:25
		Port ValueOrValueFrom `json:"port,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:Username specifies the username to use when connecting to the SMTP server.
		Username ValueOrValueFrom `json:"username,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:Password specifies the password to use when connecting to the SMTP server.
		Password ValueOrValueFrom `json:"password,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:SMTP authentication method.
		// +kubebuilder:default:Unknown
		// +kubebuilder:validation:Enum=Unknown;Plain;Login;CRAMMD5;None;OAuth2
		Auth string `json:"auth,omitempty"`

		// +kubebuilder:validation:Required
		// +kubebuilder:description:From specifies the email address to use as the sender.
		FromAddress string `json:"fromAddress,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:FromName specifies the name to use as the sender.
		FromName string `json:"fromName,omitempty"`

		// +kubebuilder:validation:Required
		// +kubebuilder:description:List of recipient e-mails.
		ToAddress []string `json:"toAddress,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:The client host name sent to the SMTP server during HELLO phase. If set to "auto" it will use the OS hostname.
		// +kubebuilder:default:auto
		ClientHost string `json:"clientHost,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:Encryption method
		// +kubebuilder:default:Auto
		// +kubebuilder:validation:Enum=Auto;None;ExplicitTLS;ImplicitTLS
		Encryption string `json:"encryption,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:Whether the message being sent is in HTML format.
		// +kubebuilder:default:false
		UseHTML bool `json:"useHTML,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:description:Whether to use the STARTTLS command (if the server supports it).
		// +kubebuilder:default:true
		UseStartTLS bool `json:"useStartTLS,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:validation:Schemaless
		// +kubebuilder:validation:Type=string
		TemplateSubject string `json:"templateSubject,omitempty"`

		// +kubebuilder:validation:Optional
		// +kubebuilder:validation:Schemaless
		// +kubebuilder:validation:Type=string
		TemplateBody string `json:"templateBody,omitempty"`
	}

	// AlertDiscordSpec defines the desired state of AlertDiscord
	AlertDiscordSpec struct {
		// +kubebuilder:validation:Required
		WebhookURL ValueOrValueFrom `json:"webhookURL"`

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
