package models

import "github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"

type (
	AlertDiscord struct {
		v1alpha1.AlertConfig
	}

	AlertEmail struct {
		v1alpha1.AlertConfig
	}
)
