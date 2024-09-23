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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kimupv1alpha1 "github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
)

// ImageReconciler reconciles a Image object
type ImageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=images,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=images/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=images/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Image object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var image kimupv1alpha1.Image

	if err := r.Client.Get(ctx, req.NamespacedName, &image); err != nil {
		log.Log.Error(err, "unable to fetch Image")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info(fmt.Sprintf("Reconciling Image %s in namespace %s", req.Name, req.Namespace))

	ok, err := image.IsChanged()
	if err != nil || ok {
		image.SetAnnotationAction(annotations.ActionRefresh)
	}

	if err := image.RefreshCheckSum(); err != nil {
		log.Log.Error(err, "unable to refresh checksum")
		return ctrl.Result{}, err
	}

	if err := r.Client.Update(ctx, &image); err != nil {
		log.Log.Error(err, "unable to update Image")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kimupv1alpha1.Image{}).
		Complete(r)
}
