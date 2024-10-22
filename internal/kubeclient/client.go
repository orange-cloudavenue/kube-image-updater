package kubeclient

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

var _ Interface = &Client{}

type (
	Client struct {
		component component
		kubernetes.Interface
		d dynamic.Interface
	}

	Interface interface {
		InterfaceKubernetes
		InterfaceKimup
	}

	InterfaceKubernetes interface {
		kubernetes.Interface
		DynamicResource(resource schema.GroupVersionResource) dynamic.NamespaceableResourceInterface
		GetPullSecretsForImage(ctx context.Context, image v1alpha1.Image) (auths K8sDockerRegistrySecretData, err error)
		GetValueOrValueFrom(ctx context.Context, namespace string, v v1alpha1.ValueOrValueFrom) (any, error)
		GetComponent() string
	}

	InterfaceKimup interface {
		Image() *ImageObj
		Alert() *AlertObj
	}

	component string
)

const (
	ComponentOperator            component = "kimup-operator"
	ComponentController          component = "kimup-controller"
	ComponentAdmissionController component = "kimup-admission-controller"
)

func init() {
	if flag.Lookup("kubeconfig") == nil {
		flag.String("kubeconfig", "", "path to the kubeconfig file")
	}
}

// New creates a new kubernetes client
// kubeConfigPath is the path to the kubeconfig file (empty for in-cluster)
func New(kubeConfigPath string, c component) (Interface, error) {
	config, err := getConfig(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return NewFromRestConfig(config, c)
}

// NewFromRestConfig creates a new kubernetes client from a rest config
func NewFromRestConfig(config *rest.Config, c component) (*Client, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		Interface: client,
		d:         dynamicClient,
		component: c,
	}, nil
}

func getConfig(kubeConfigPath string) (config *rest.Config, err error) {
	if kubeConfigPath != "" {
		// use the current context in kubeconfig
		return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	}

	// creates the in-cluster config
	return rest.InClusterConfig()
}

// DynamicResource returns a dynamic resource
func (c *Client) DynamicResource(resource schema.GroupVersionResource) dynamic.NamespaceableResourceInterface {
	return c.d.Resource(resource)
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
		secret, err := c.CoreV1().Secrets(image.Namespace).Get(ctx, ip.Name, metav1.GetOptions{})
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

			auths.Auths[k] = v
		}
	}

	return auths, nil
}

func (c *Client) GetComponent() string {
	return string(c.component)
}
