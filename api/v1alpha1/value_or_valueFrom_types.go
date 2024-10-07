package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type (
	ValueOrValueFrom struct {
		// Value is a string value to assign to the key.
		// if ValueFrom is specified, this value is ignored.
		// +optional
		Value string `json:"value,omitempty"`

		// ValueFrom is a reference to a field in a secret or config map.
		// +optional
		ValueFrom *ValueFromSource `json:"valueFrom,omitempty"`
	}

	// ValueFromSource is a reference to a field in a secret or config map.
	ValueFromSource struct {
		// SecretKeyRef is a reference to a field in a secret.
		// +optional
		SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`

		// ConfigMapKeyRef is a reference to a field in a config map.
		// +optional
		ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`

		// AlertConfigRef is a reference to a field in an alert configuration.
		// +optional
		AlertConfigRef *corev1.LocalObjectReference `json:"alertConfigRef,omitempty"`
	}
)
