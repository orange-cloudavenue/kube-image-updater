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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

// ImageReconciler reconciles a Image object
type AlertDiscordReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=alertdiscord,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=alertdiscord/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=alertdiscord/finalizers,verbs=update

func (r *AlertDiscordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var aDiscord v1alpha1.AlertDiscord

	if err := r.Client.Get(ctx, req.NamespacedName, &aDiscord); err != nil {
		log.Log.Error(err, "unable to fetch AlertDiscord")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// * Status

	if err := r.Status().Update(ctx, &aDiscord); err != nil {
		log.Log.Error(err, "unable to update AlertDiscord status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertDiscordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AlertDiscord{}).
		Complete(r)
}
