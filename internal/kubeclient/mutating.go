package kubeclient

import (
	"context"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

type (
	AdmissionControllerObj struct {
		InterfaceKubernetes
	}
)

// AdmissionController returns an AdmissionController object
func (c *Client) AdmissionController() *AdmissionControllerObj {
	return NewAdmissionController(c)
}

func NewAdmissionController(k InterfaceKubernetes) *AdmissionControllerObj {
	return &AdmissionControllerObj{
		InterfaceKubernetes: k,
	}
}

func (a *AdmissionControllerObj) GetMutatingConfiguration(ctx context.Context, name string) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	return a.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(ctx, name, metav1.GetOptions{})
}

func (a *AdmissionControllerObj) CreateOrUpdateMutatingConfiguration(ctx context.Context, name string, svc admissionregistrationv1.ServiceReference, policy admissionregistrationv1.FailurePolicyType) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	mutatingWebhookConfig := a.buildMutatingConfiguration(name, svc, policy)
	if _, err := a.GetMutatingConfiguration(ctx, name); err != nil {
		if apierrors.IsNotFound(err) {
			return a.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(ctx, mutatingWebhookConfig, metav1.CreateOptions{})
		}
		return nil, err
	}

	return a.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(ctx, mutatingWebhookConfig, metav1.UpdateOptions{})
}

func (a *AdmissionControllerObj) buildMutatingConfiguration(name string, svc admissionregistrationv1.ServiceReference, policy admissionregistrationv1.FailurePolicyType) *admissionregistrationv1.MutatingWebhookConfiguration {
	return &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"cert-manager.io/inject-ca-from": "kimup-operator/kimup-webhook-serving-cert",
			},
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{{
			Name:                    "image-tag.kimup.cloudavenue.io",
			AdmissionReviewVersions: []string{"v1", "v1beta1"},
			SideEffects:             utils.ToPTR(admissionregistrationv1.SideEffectClassNone),
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				Service: &svc,
			},
			Rules: []admissionregistrationv1.RuleWithOperations{
				{
					Operations: []admissionregistrationv1.OperationType{
						admissionregistrationv1.Update,
						admissionregistrationv1.Create,
					},
					Rule: admissionregistrationv1.Rule{
						APIGroups:   []string{"*"},
						APIVersions: []string{"v1"},
						Resources:   []string{"pods"},
						Scope:       utils.ToPTR(admissionregistrationv1.NamespacedScope),
					},
				},
			},
			NamespaceSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{{
					Key:      string(annotations.KeyEnabled),
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{"true", "yes"},
				}},
			},
			FailurePolicy: utils.ToPTR(policy),
		}},
	}
}
