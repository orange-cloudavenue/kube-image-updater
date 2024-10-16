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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kimupv1alpha1 "github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

// ImageReconciler reconciles a Image object
type ImageReconciler struct {
	client.Client
	KubeAPIClient *kubeclient.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
}

type ImageEvent string

const (
	ImageUpdate ImageEvent = "ImageUpdate"
)

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
	xlog := log.WithContext(ctx).WithFields(logrus.Fields{
		"namespace": req.Namespace,
		"name":      req.Name,
	})

	var image kimupv1alpha1.Image

	if err := r.Client.Get(ctx, req.NamespacedName, &image); err != nil {
		if client.IgnoreNotFound(err) != nil {
			xlog.WithError(err).Error("could not get the image object")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	an := annotations.New(ctx, &image)

	xlog.
		WithFields(logrus.Fields{
			"name":      req.Name,
			"namespace": req.Namespace,
		}).Info("Reconciling Image")
	equal, err := an.CheckSum().IsEqual(image.Spec)
	if err != nil || !equal {
		an.Action().Set(annotations.ActionReload)
		r.Recorder.Event(&image, "Normal", string(ImageUpdate), "Image configuration has changed. Reloading image.")
	}

	if an.CheckSum().IsNull() || !equal {
		if err := an.CheckSum().Set(image.Spec); err != nil {
			xlog.WithError(err).Error("unable to set checksum")
			return ctrl.Result{}, err
		}
		if err := r.Client.Update(ctx, &image); err != nil {
			xlog.WithError(err).Error("unable to update Image")
			return ctrl.Result{}, err
		}
	}

	// * Status

	image.SetStatusTag(an.Tag().Get())
	if err := r.Status().Update(ctx, &image); err != nil {
		xlog.WithError(err).Error("unable to update Image status")
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
