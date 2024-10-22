package fakekubeclient

import (
	"context"

	"github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	dFake "k8s.io/client-go/dynamic/fake"
	kFake "k8s.io/client-go/kubernetes/fake"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
)

var _ kubeclient.Interface = &FakeKubeClient{}

type FakeKubeClient struct {
	mock.Mock
	kubeclient.InterfaceKubernetes
}

func NewFakeKubeClient() *FakeKubeClient {
	return &FakeKubeClient{
		InterfaceKubernetes: &kubeclient.Client{
			Interface: kFake.NewSimpleClientset(),
		},
	}
}

func (f *FakeKubeClient) DynamicResource(resource schema.GroupVersionResource) dynamic.NamespaceableResourceInterface {
	return dFake.NewSimpleDynamicClient(runtime.NewScheme()).Resource(resource)
}

func (f *FakeKubeClient) GetPullSecretsForImage(ctx context.Context, image v1alpha1.Image) (auths kubeclient.K8sDockerRegistrySecretData, err error) {
	args := f.Called(ctx, image)
	return args.Get(0).(kubeclient.K8sDockerRegistrySecretData), args.Error(1)
}

func (f *FakeKubeClient) GetValueOrValueFrom(ctx context.Context, namespace string, v v1alpha1.ValueOrValueFrom) (any, error) {
	args := f.Called(ctx, namespace, v)
	return args.Get(0), args.Error(1)
}

func (f *FakeKubeClient) GetComponent() string {
	return "unit-test"
}

func (f *FakeKubeClient) Image() *kubeclient.ImageObj {
	return kubeclient.NewImage(f)
}

func (f *FakeKubeClient) Alert() *kubeclient.AlertObj {
	return kubeclient.NewAlert(f)
}

func (f *FakeKubeClient) CreateFakeImage(image v1alpha1.Image) error {
	u, err := kubeclient.EncodeUnstructured(image)
	if err != nil {
		return err
	}

	_, err = f.DynamicResource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: "images",
	}).
		Namespace(image.Namespace).
		Create(
			context.Background(),
			u,
			v1.CreateOptions{},
		)
	return err
}
