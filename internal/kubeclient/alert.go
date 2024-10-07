package kubeclient

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

type (
	AlertObj struct {
		c           *Client
		alertClient dynamic.NamespaceableResourceInterface
	}
)

// Alert() returns an alert object
func (c *Client) Alert() *AlertObj {
	return &AlertObj{
		c: c,
		alertClient: c.d.Resource(schema.GroupVersionResource{
			Group:    v1alpha1.GroupVersion.Group,
			Version:  v1alpha1.GroupVersion.Version,
			Resource: "alertconfig",
		}),
	}
}

// Get retrieves an Alert object by its name within the specified namespace.
// It takes a context, the namespace, and the name of the Alert as parameters.
// If the Alert is found, it returns a pointer to the Alert object and a nil error.
// If there is an error during the retrieval process, it returns nil and the error encountered.
func (a *AlertObj) Get(ctx context.Context, name string) (v1alpha1.AlertConfig, error) {
	u, err := a.alertClient.Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return v1alpha1.AlertConfig{}, err
	}

	return decodeUnstructured[v1alpha1.AlertConfig](u)
}

// List retrieves a list of AlertObj instances from the specified namespace.
// It takes a context, the namespace as a string, and list options.
// Returns a pointer to a List of AlertObj and an error if the operation fails.
func (a *AlertObj) List(ctx context.Context, opts v1.ListOptions) (v1alpha1.AlertConfigList, error) {
	u, err := a.alertClient.List(ctx, opts)
	if err != nil {
		return v1alpha1.AlertConfigList{}, err
	}

	return decodeUnstructured[v1alpha1.AlertConfigList](u)
}

// Update updates an existing alert in the specified namespace.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for cancellation and deadlines.
//   - namespace: The namespace where the alert is located.
//   - alert: A pointer to the alert object that needs to be updated.
//
// Returns:
//   - An error if the update operation fails; otherwise, it returns nil.
func (a *AlertObj) Update(ctx context.Context, alert v1alpha1.AlertConfig) error {
	u, err := encodeUnstructured(alert)
	if err != nil {
		return err
	}

	_, err = a.alertClient.Update(ctx, u, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
