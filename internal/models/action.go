package models

import (
	"context"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
)

type (
	ActionInterface interface {
		Init(kubeClient *kubeclient.Client, tags Tags, image *v1alpha1.Image, data v1alpha1.ValueOrValueFrom)
		Execute(context.Context) error
		GetName() ActionName
		GetActualTag() string
		GetNewTag() string
		GetAvailableTags() []string
	}

	ActionName string
)

// String returns the string representation of the action name.
func (n ActionName) String() string {
	return string(n)
}
