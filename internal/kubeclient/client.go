package kubeclient

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

type (
	Client struct {
		c *kubernetes.Clientset
		d *dynamic.DynamicClient
	}
)

func init() {
	flag.String("kubeconfig", "", "path to the kubeconfig file")
}

// New creates a new kubernetes client
// kubeConfigPath is the path to the kubeconfig file (empty for in-cluster)
func New(kubeConfigPath string) (*Client, error) {
	client, dynamicClient, err := newClientK8s(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return &Client{c: client, d: dynamicClient}, nil
}

func getConfig(kubeConfigPath string) (config *rest.Config, err error) {
	if kubeConfigPath != "" {
		// use the current context in kubeconfig
		return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	}

	// creates the in-cluster config
	return rest.InClusterConfig()
}

func newClientK8s(kubeConfigPath string) (*kubernetes.Clientset, *dynamic.DynamicClient, error) {
	config, err := getConfig(kubeConfigPath)
	if err != nil {
		return nil, nil, err
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	d, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	return c, d, nil
}

// GetKubeClient returns the standard kubernetes client
func (c *Client) GetKubeClient() *kubernetes.Clientset {
	return c.c
}

// GetDynamicClient returns the dynamic kubernetes client
func (c *Client) GetDynamicClient() *dynamic.DynamicClient {
	return c.d
}

// ! Images

func (c *Client) cImage() dynamic.NamespaceableResourceInterface {
	return c.d.Resource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: "images",
	})
}

func (c *Client) listImages(ctx context.Context, namespace string) (list v1alpha1.ImageList, err error) {
	var v *unstructured.UnstructuredList

	if namespace == "" {
		v, err = c.cImage().List(ctx, metav1.ListOptions{})
	} else {
		v, err = c.cImage().Namespace(namespace).List(ctx, metav1.ListOptions{})
	}
	if err != nil {
		return list, fmt.Errorf("failed to list resources: %w", err)
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(v.UnstructuredContent(), &list); err != nil {
		return list, fmt.Errorf("failed to convert resource: %w", err)
	}

	return
}

// ListAllImages lists all images in all namespaces
func (c *Client) ListAllImages(ctx context.Context) (list v1alpha1.ImageList, err error) {
	return c.listImages(ctx, "")
}

// ListImages lists all images in a namespace
func (c *Client) ListImages(ctx context.Context, namespace string) (list v1alpha1.ImageList, err error) {
	return c.listImages(ctx, namespace)
}

// GetImage gets an image in a namespace
func (c *Client) GetImage(ctx context.Context, namespace, name string) (image v1alpha1.Image, err error) {
	if namespace == "" {
		return image, fmt.Errorf("namespace is required")
	}

	if name == "" {
		return image, fmt.Errorf("name is required")
	}

	v, err := c.cImage().Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return image, fmt.Errorf("failed to get resource: %w", err)
	}

	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(v.UnstructuredContent(), &image); err != nil {
		return image, fmt.Errorf("failed to convert resource: %w", err)
	}

	return
}

// SetImage sets an image in a namespace
func (c *Client) SetImage(ctx context.Context, image v1alpha1.Image) (err error) {
	unstructedImage, err := runtime.DefaultUnstructuredConverter.
		ToUnstructured(&image)
	if err != nil {
		return fmt.Errorf("failed to convert resource: %w", err)
	}

	_, err = c.cImage().Namespace(image.Namespace).Update(ctx, &unstructured.Unstructured{Object: unstructedImage}, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	return
}

// FindImage finds an image in a namespace
func (c *Client) FindImage(ctx context.Context, namespace, name string) (image v1alpha1.Image, err error) {
	l, err := c.listImages(ctx, namespace)
	if err != nil {
		return image, fmt.Errorf("failed to list images: %w", err)
	}
	for _, i := range l.Items {
		if i.GetImageWithoutTag() == utils.ImageParser(name).GetImageWithoutTag() {
			return i, nil
		}
	}

	return image, fmt.Errorf("image not found")
}

type K8sDockerRegistrySecretData struct {
	Auths map[string]K8sDockerRegistrySecret `json:"auths"`
}

type K8sDockerRegistrySecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth"`
}

func (c *Client) GetPullSecretsForImage(ctx context.Context, image v1alpha1.Image) (auths K8sDockerRegistrySecretData, err error) {
	auths.Auths = make(map[string]K8sDockerRegistrySecret)

	for _, ip := range image.Spec.ImagePullSecrets {
		secret, err := c.GetKubeClient().CoreV1().Secrets(image.Namespace).Get(ctx, ip.Name, metav1.GetOptions{})
		if err != nil {
			continue
		}

		if secret.Type != v1.SecretTypeDockerConfigJson {
			continue
		}

		auth := K8sDockerRegistrySecretData{}
		if err := json.Unmarshal(secret.Data[v1.DockerConfigJsonKey], &auth); err != nil {
			return auths, fmt.Errorf("failed to unmarshal secret: %w", err)
		}

		for k, v := range auth.Auths {
			if v.Username == "" || v.Password == "" {
				continue
			}

			for _, i := range []string{"https://", "http://"} {
				k = strings.TrimPrefix(k, i)
			}

			log.Debugf("Found auth for %s", k)
			auths.Auths[k] = v
		}
	}

	return auths, nil
}
