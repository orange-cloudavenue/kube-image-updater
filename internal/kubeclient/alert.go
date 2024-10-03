package kubeclient

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

type (
	alert struct {
		c *Client
	}

	AlertList struct {
		alert
	}

	AlertObj[Obj, List any] struct {
		alert
		alertClient dynamic.NamespaceableResourceInterface
	}
)

// Alert() returns an alert object
func (c *Client) Alert() *AlertList {
	return &AlertList{alert: alert{
		c: c,
	}}
}

// Discord returns an alert object for Discord
func (a *AlertList) Discord() *AlertObj[v1alpha1.AlertDiscord, v1alpha1.AlertDiscordList] {
	return &AlertObj[v1alpha1.AlertDiscord, v1alpha1.AlertDiscordList]{
		alert: a.alert,
		alertClient: a.c.d.Resource(schema.GroupVersionResource{
			Group:    v1alpha1.GroupVersion.Group,
			Version:  v1alpha1.GroupVersion.Version,
			Resource: "alertdiscords",
		}),
	}
}

// Get retrieves an Alert object by its name within the specified namespace.
// It takes a context, the namespace, and the name of the Alert as parameters.
// If the Alert is found, it returns a pointer to the Alert object and a nil error.
// If there is an error during the retrieval process, it returns nil and the error encountered.
func (a *AlertObj[Obj, List]) Get(ctx context.Context, namespace, name string) (Obj, error) {
	u, err := a.alertClient.Namespace(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return *new(Obj), err
	}

	return decodeUnstructured[Obj](u)
}

// List retrieves a list of AlertObj instances from the specified namespace.
// It takes a context, the namespace as a string, and list options.
// Returns a pointer to a List of AlertObj and an error if the operation fails.
func (a *AlertObj[Obj, List]) List(ctx context.Context, namespace string, opts v1.ListOptions) (List, error) {
	u, err := a.alertClient.Namespace(namespace).List(ctx, opts)
	if err != nil {
		return *new(List), err
	}

	return decodeUnstructured[List](u)
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
func (a *AlertObj[Obj, List]) Update(ctx context.Context, namespace string, alert Obj) error {
	u, err := encodeUnstructured(alert)
	if err != nil {
		return err
	}

	_, err = a.alertClient.Namespace(namespace).Update(ctx, u, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
