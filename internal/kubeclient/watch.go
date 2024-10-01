package kubeclient

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

type (
	WatchInterface[T any] struct {
		Type  watch.EventType
		Value *T
	}
)

// WatchEventsImage watches events for an image in a namespace
func (c *Client) WatchEventsImage(ctx context.Context) (chan WatchInterface[v1alpha1.Image], error) {
	x, err := c.cImage().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	ch := make(chan WatchInterface[v1alpha1.Image])

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				x.Stop()
				return
			case event, ok := <-x.ResultChan():
				if !ok {
					close(ch)
					return
				}

				image, err := decodeUnstructured[*v1alpha1.Image](event.Object.(*unstructured.Unstructured))
				if err != nil {
					log.Errorf("Error decoding image: %v", err)
					continue
				}

				ch <- WatchInterface[v1alpha1.Image]{Type: event.Type, Value: image}
			}
		}
	}()

	return ch, nil
}
