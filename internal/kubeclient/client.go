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
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

type (
	Client struct {
		c *kubernetes.Clientset
		d *dynamic.DynamicClient
	}
)

func init() {
	if flag.Lookup("kubeconfig") == nil {
		flag.String("kubeconfig", "", "path to the kubeconfig file")
	}
}

// New creates a new kubernetes client
// kubeConfigPath is the path to the kubeconfig file (empty for in-cluster)
func New(kubeConfigPath string) (*Client, error) {
	config, err := getConfig(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return NewFromRestConfig(config)
}

// NewFromRestConfig creates a new kubernetes client from a rest config
func NewFromRestConfig(config *rest.Config) (*Client, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
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

// GetKubeClient returns the standard kubernetes client
func (c *Client) GetKubeClient() *kubernetes.Clientset {
	return c.c
}

// GetDynamicClient returns the dynamic kubernetes client
func (c *Client) GetDynamicClient() *dynamic.DynamicClient {
	return c.d
}

type UnstructuredFunc interface {
	UnstructuredContent() map[string]interface{}
}

func decodeUnstructured[T any](v UnstructuredFunc) (t T, err error) {
	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(v.UnstructuredContent(), &t); err != nil {
		return t, fmt.Errorf("failed to convert resource: %w", err)
	}

	return
}

func encodeUnstructured[T any](t T) (*unstructured.Unstructured, error) {
	x, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&t)
	if err != nil {
		return nil, fmt.Errorf("failed to convert resource: %w", err)
	}

	return &unstructured.Unstructured{Object: x}, nil
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
