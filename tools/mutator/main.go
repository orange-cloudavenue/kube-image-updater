// This tool is used to create mutating configuration for the admission controller webhook.

package main

import (
	"context"
	"flag"
	"os"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

func main() {
	kubeconfig := flag.Lookup("kubeconfig").Value.String()

	if kubeconfig == "" {
		// Get home directory
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		kubeconfig = home + "/.kube/config"
	}

	// kubernetes golang library provide flag "kubeconfig" to specify the path to the kubeconfig file
	k, err := kubeclient.New(kubeconfig, kubeclient.ComponentOperator)
	if err != nil {
		log.WithError(err).Panic("Error creating kubeclient")
	}

	_, err = k.Mutator().CreateOrUpdateMutatingConfiguration(
		context.Background(),
		models.MutatorMutatingWebhookConfigurationName,
		admissionregistrationv1.ServiceReference{
			Name:      "mutator",
			Namespace: "kimup-operator",
			Path:      &models.MutatorWebhookPathMutateImageTag,
		},
		admissionregistrationv1.Fail,
	)
	if err != nil {
		log.WithError(err).Panic("Error creating or updating mutating configuration")
	}
}
