package kubeclient

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

type (
	ImageObj struct {
		InterfaceKubernetes
		imageClient dynamic.NamespaceableResourceInterface
	}
)

// Image returns an image object
func (c *Client) Image() *ImageObj {
	return NewImage(c)
}

func NewImage(k InterfaceKubernetes) *ImageObj {
	return &ImageObj{
		InterfaceKubernetes: k,
		imageClient: k.DynamicResource(schema.GroupVersionResource{
			Group:    v1alpha1.GroupVersion.Group,
			Version:  v1alpha1.GroupVersion.Version,
			Resource: "images",
		}),
	}
}

// Get retrieves an Image object by its name within the specified namespace.
// It takes a context, the namespace, and the name of the Image as parameters.
// If the Image is found, it returns a pointer to the Image object and a nil error.
// If there is an error during the retrieval process, it returns nil and the error encountered.
func (i *ImageObj) Get(ctx context.Context, namespace, name string) (v1alpha1.Image, error) {
	u, err := i.imageClient.Namespace(namespace).Get(ctx, name, v1.GetOptions{})
	if err != nil {
		return v1alpha1.Image{}, err
	}

	return decodeUnstructured[v1alpha1.Image](u)
}

// List retrieves a list of images from the specified namespace.
// It takes a context, the namespace as a string, and list options.
// Returns a pointer to a List of images and an error if the operation fails.
func (i *ImageObj) List(ctx context.Context, namespace string, opts v1.ListOptions) (v1alpha1.ImageList, error) {
	return i.listImages(ctx, namespace, opts)
}

// ListAll retrieves a list of images from all namespaces.
// It takes a context and list options as parameters.
// Returns a pointer to a List of images and an error if the operation fails.
func (i *ImageObj) ListAll(ctx context.Context, opts v1.ListOptions) (v1alpha1.ImageList, error) {
	return i.listImages(ctx, "", opts)
}

// listImages lists all images
// It takes a context and a namespace as parameters.
// if namespace is empty, it lists all images in all namespaces.
// Returns a pointer to a List of images and an error if the operation fails.
func (i *ImageObj) listImages(ctx context.Context, namespace string, opts v1.ListOptions) (v1alpha1.ImageList, error) {
	var (
		err error
		u   *unstructured.UnstructuredList
	)

	if namespace == "" {
		u, err = i.imageClient.List(ctx, opts)
	} else {
		u, err = i.imageClient.Namespace(namespace).List(ctx, opts)
	}
	if err != nil {
		return v1alpha1.ImageList{}, fmt.Errorf("failed to list resources: %w", err)
	}

	return decodeUnstructured[v1alpha1.ImageList](u)
}

// Update the image object in the specified namespace.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for cancellation and deadlines.
//   - namespace: The namespace in which the image object resides.
//   - image: A pointer to the image object to be updated.
//
// Returns:
//   - An error if the update operation fails; otherwise, it returns nil.
func (i *ImageObj) Update(ctx context.Context, image v1alpha1.Image) error {
	u, err := encodeUnstructured(image)
	if err != nil {
		return err
	}

	_, err = i.imageClient.Namespace(image.Namespace).Update(ctx, u, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return err
}

// Find finds an image by its image name. Example: `docker.io/library/nginx:latest`
// It takes a context and the image name as parameters.
// Returns a pointer to the Image object and an error if the operation fails.
func (i *ImageObj) Find(ctx context.Context, namespace, imageName string) (v1alpha1.Image, error) {
	images, err := i.listImages(ctx, namespace, v1.ListOptions{})
	if err != nil {
		return v1alpha1.Image{}, err
	}

	for _, image := range images.Items {
		if image.Spec.Image == imageName {
			return image, nil
		}
	}

	return v1alpha1.Image{}, fmt.Errorf("image %s %w", imageName, ErrNotFound)
}

// Watch watches for changes to the image object.
// It takes a context and the namespace as parameters.
// Returns a channel of WatchInterface[v1alpha1.Image] and an error if the operation fails.
func (i *ImageObj) Watch(ctx context.Context) (chan WatchInterface[v1alpha1.Image], error) {
	x, err := i.imageClient.Watch(ctx, v1.ListOptions{})
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

				image, err := decodeUnstructured[v1alpha1.Image](event.Object.(*unstructured.Unstructured))
				if err != nil {
					log.WithError(err).Error("Failed to decode image")
					continue
				}

				ch <- WatchInterface[v1alpha1.Image]{Type: event.Type, Value: image}
			}
		}
	}()

	return ch, nil
}
