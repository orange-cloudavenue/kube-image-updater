/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"github.com/sirupsen/logrus"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	KubeAPIClient *kubeclient.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
}

// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Image object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	xlog := log.WithContext(ctx).WithFields(logrus.Fields{
		"namespace": req.Namespace,
		"name":      req.Name,
	})

	xlog.Info("Reconciling Namespace")

	var (
		ns                   corev1.Namespace
		foundInMutatorConfig bool
	)

	if err := r.Client.Get(ctx, req.NamespacedName, &ns); err != nil {
		if client.IgnoreNotFound(err) == nil {
			xlog.WithError(err).Error("could not get the image object")
			foundInMutatorConfig = true // Force rebuilding the mutating configuration
		} else {
			return ctrl.Result{}, err
		}
	}

	// get mutator configuration
	mutator, _ := r.KubeAPIClient.Mutator().GetMutatingConfiguration(ctx, models.MutatorWebhookConfigurationName)
	// ignore error, we will create it if it does not exist
	if mutator != nil {
		wName := kubeclient.NamespaceMatchConditionBuilder{}.New(req.Name).GetName()
		for _, webhook := range mutator.Webhooks {
			if webhook.Name == wName {
				foundInMutatorConfig = true
				break
			}
		}
	}

	an := annotations.New(ctx, &ns)

	if an.Enabled().Get() || foundInMutatorConfig {
		_, err := r.KubeAPIClient.Mutator().CreateOrUpdateMutatingConfiguration(
			ctx,
			models.MutatorWebhookConfigurationName,
			admissionregistrationv1.ServiceReference{
				Name:      "mutator",
				Namespace: "kimup-operator",
				Path:      &models.MutatorWebhookPathMutateImageTag,
			},
			admissionregistrationv1.Fail,
		)
		if err != nil {
			xlog.WithError(err).Error("could not create or update mutating configuration")
			return ctrl.Result{RequeueAfter: utils.RandomSecondInRange(1, 7)}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(r)
}
