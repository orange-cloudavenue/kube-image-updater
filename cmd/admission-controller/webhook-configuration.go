package main

import (
	"bytes"
	"context"
	"reflect"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	client "github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

// createOrUpdateMutatingWebhookConfiguration creates or updates the mutating webhook configuration
// for the webhook service. The CA is generated and used for the webhook.
// This function create the request to the Kubernetes API server to create or update the mutating webhook configuration.
func createOrUpdateMutatingWebhookConfiguration(caPEM *bytes.Buffer, webhookService, webhookNamespace string, k client.Interface) error {
	mutatingWebhookConfigV1Client := k.AdmissionregistrationV1()

	var clientConfig admissionregistrationv1.WebhookClientConfig
	switch insideCluster {
	case true:
		clientConfig = admissionregistrationv1.WebhookClientConfig{
			Service: &admissionregistrationv1.ServiceReference{
				Name:      webhookService,
				Namespace: webhookNamespace,
				Path:      &webhookPathMutate,
			},
		}
	case false:
		// the webhook is running outside the cluster
		// Please note that the webhook service must be accessible from the Kubernetes cluster.
		// Each time you change webhook service name, namespace, or port, you need to update the MutatingWebhookConfiguration
		// Also you need to modifiy the manifest MutatingWebhookConfiguration to push new caPEM to allow client to trust the webhook
		// The caPEM is generated and printed to the logs when the webhook starts for outside cluster
		url := "https://" + webhookService + "." + webhookNamespace + webhookPort + webhookPathMutate
		clientConfig = admissionregistrationv1.WebhookClientConfig{
			CABundle: caPEM.Bytes(),
			URL:      &url,
		}
	}
	log.Debug("Creating or updating the mutatingwebhookconfiguration")
	fail := admissionregistrationv1.Fail
	sideEffect := admissionregistrationv1.SideEffectClassNone
	mutatingWebhookConfig := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: webhookConfigName,
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{{
			Name:                    webhookService + "." + webhookNamespace,
			AdmissionReviewVersions: []string{"v1", "v1beta1"},
			SideEffects:             &sideEffect,
			ClientConfig:            clientConfig,

			Rules: []admissionregistrationv1.RuleWithOperations{
				{
					Operations: []admissionregistrationv1.OperationType{
						admissionregistrationv1.Update,
						admissionregistrationv1.Create,
					},
					Rule: admissionregistrationv1.Rule{
						APIGroups:   []string{""},
						APIVersions: []string{"v1"},
						Resources:   []string{"pods"},
						// TODO - add namespace scope
						// Scope:       "*",
					},
				},
			},
			FailurePolicy: &fail,
		}},
	}

	// check if the mutatingwebhookconfiguration already exists
	foundWebhookConfig, err := mutatingWebhookConfigV1Client.MutatingWebhookConfigurations().Get(context.TODO(), webhookConfigName, metav1.GetOptions{})
	switch {
	case err != nil && apierrors.IsNotFound(err):
		if _, err := mutatingWebhookConfigV1Client.MutatingWebhookConfigurations().Create(context.TODO(), mutatingWebhookConfig, metav1.CreateOptions{}); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		// there is an existing mutatingWebhookConfiguration
		if len(foundWebhookConfig.Webhooks) != len(mutatingWebhookConfig.Webhooks) ||
			!(foundWebhookConfig.Webhooks[0].Name == mutatingWebhookConfig.Webhooks[0].Name &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].AdmissionReviewVersions, mutatingWebhookConfig.Webhooks[0].AdmissionReviewVersions) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].SideEffects, mutatingWebhookConfig.Webhooks[0].SideEffects) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].FailurePolicy, mutatingWebhookConfig.Webhooks[0].FailurePolicy) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].Rules, mutatingWebhookConfig.Webhooks[0].Rules) &&
				// reflect.DeepEqual(foundWebhookConfig.Webhooks[0].NamespaceSelector, mutatingWebhookConfig.Webhooks[0].NamespaceSelector) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.CABundle, mutatingWebhookConfig.Webhooks[0].ClientConfig.CABundle) &&
				// reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.Service, mutatingWebhookConfig.Webhooks[0].ClientConfig.Service) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.URL, mutatingWebhookConfig.Webhooks[0].ClientConfig.URL)) {
			mutatingWebhookConfig.ObjectMeta.ResourceVersion = foundWebhookConfig.ObjectMeta.ResourceVersion
			if _, err := mutatingWebhookConfigV1Client.MutatingWebhookConfigurations().Update(context.TODO(), mutatingWebhookConfig, metav1.UpdateOptions{}); err != nil {
				return err
			}
		}
	}

	return nil
}
