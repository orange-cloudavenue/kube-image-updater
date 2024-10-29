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
	MutatorObj struct {
		InterfaceKubernetes
	}
)

// Mutator returns an Mutator object
func (c *Client) Mutator() *MutatorObj {
	return NewMutator(c)
}

func NewMutator(k InterfaceKubernetes) *MutatorObj {
	return &MutatorObj{
		InterfaceKubernetes: k,
	}
}

func (a *MutatorObj) GetMutatingConfiguration(ctx context.Context, name string) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	return a.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(ctx, name, metav1.GetOptions{})
}

func (a *MutatorObj) CreateOrUpdateMutatingConfiguration(ctx context.Context, name string, svc admissionregistrationv1.ServiceReference, policy admissionregistrationv1.FailurePolicyType) (*admissionregistrationv1.MutatingWebhookConfiguration, error) {
	// Get All Namespaces with the "enabled" label
	nsList, err := a.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	mwc, err := a.GetMutatingConfiguration(ctx, name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			mwc = &admissionregistrationv1.MutatingWebhookConfiguration{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
					Annotations: map[string]string{
						"cert-manager.io/inject-ca-from": "kimup-operator/kimup-webhook-serving-cert",
					},
				},
			}
		}
	}

	// reset webhooks settings
	mwc.Webhooks = []admissionregistrationv1.MutatingWebhook{}

	for _, ns := range nsList.Items {
		nsAnnotation := annotations.New(ctx, &ns)
		if !nsAnnotation.Enabled().Get() {
			continue
		}

		mwc.Webhooks = append(mwc.Webhooks, a.buildMutatingWebhookConfiguration(svc, policy, &namespaceMatchConditionBuilder{namespace: ns.Name}))
	}

	// Add the default matchCondition (All pods with annotation enabled == true)
	mwc.Webhooks = append(mwc.Webhooks, a.buildMutatingWebhookConfiguration(svc, policy, &defaultMatchConditionBuilder{}))

	if mwc.UID == "" {
		return a.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(ctx, mwc, metav1.CreateOptions{})
	}

	return a.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(ctx, mwc, metav1.UpdateOptions{})
}

func (a *MutatorObj) buildMutatingWebhookConfiguration(svc admissionregistrationv1.ServiceReference, policy admissionregistrationv1.FailurePolicyType, matchConditionBuilder matchConditionBuilderInterface) admissionregistrationv1.MutatingWebhook {
	return admissionregistrationv1.MutatingWebhook{
		Name:                    matchConditionBuilder.getName() + ".image-tag.kimup.cloudavenue.io",
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
		MatchConditions: matchConditionBuilder.buildMatchCondition(),
		FailurePolicy:   utils.ToPTR(policy),
	}
}