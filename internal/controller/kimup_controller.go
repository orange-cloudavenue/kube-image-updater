package controller

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

// KimupReconciler reconciles a Image object
type KimupReconciler struct {
	client.Client
	KubeAPIClient *kubeclient.Client
	Scheme        *runtime.Scheme
	Recorder      record.EventRecorder
}

//+kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=kimups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=kimups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kimup.cloudavenue.io,resources=kimups/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *KimupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) { //nolint:gocyclo
	xlog := log.WithContext(ctx).WithFields(logrus.Fields{
		"namespace": req.Namespace,
		"name":      req.Name,
	})

	var kim v1alpha1.Kimup
	if err := r.Get(ctx, req.NamespacedName, &kim); err != nil {
		if client.IgnoreNotFound(err) != nil {
			xlog.WithError(err).Error("could not get the kimup object")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var (
		rescheduleAfter   = false
		resourcesToCreate = make([]Object, 0)
		resourcesToUpdate = make([]Object, 0)
	)

	if kim.Spec.Controller != nil {
		switch {
		case kim.Status.Controller.State == "":
			resourcesToCreate = append(resourcesToCreate, GetKimupControllerResources(ctx, kim)...)
		default:
			resourcesToUpdate = append(resourcesToUpdate, GetKimupControllerResources(ctx, kim)...)
		}
	}

	if kim.Spec.AdmissionController != nil {
		switch {
		case kim.Status.AdmissionController.State == "":
			resourcesToCreate = append(resourcesToCreate, GetKimupAdmissionResources(ctx, kim)...)
		default:
			resourcesToUpdate = append(resourcesToUpdate, GetKimupAdmissionResources(ctx, kim)...)
		}
	}

	for _, resource := range resourcesToCreate {
		if err := r.Create(ctx, resource.obj); err != nil {
			if client.IgnoreAlreadyExists(err) == nil {
				resourcesToUpdate = append(resourcesToUpdate, resource)
				continue
			}
			xlog.WithError(err).Error("could not create the resource")
			rescheduleAfter = true
			continue
		}
	}

	for _, resource := range resourcesToUpdate {
		if err := r.Update(ctx, resource.obj); err != nil {
			xlog.WithError(err).Error("could not update the resource")
			rescheduleAfter = true
			continue
		}
	}

	time.Sleep(1 * time.Second)

	// Get all the resources for the kimup object and check if they are ready
	// If they are not ready, requeue the request
	// If they are ready, update the status of the kimup object

	allResources := []Object{}
	allResources = append(allResources, resourcesToCreate...)
	allResources = append(allResources, resourcesToUpdate...)

	for _, resource := range allResources {
		switch resource.kind {
		case "Deployment":
			var deployment appsv1.Deployment
			if err := r.Get(ctx, client.ObjectKeyFromObject(resource.obj), &deployment); err != nil {
				xlog.
					WithError(err).
					WithFields(logrus.Fields{
						"namespace": deployment.Namespace,
						"name":      deployment.Name,
						"kind":      deployment.Kind,
					}).Error("could not get the deployment")
				rescheduleAfter = true
				continue
			}

			switch deployment.Name {
			case KimupControllerName:
				kim.Status.Controller.State = StateResourcesCreated
			case KimupAdmissionControllerName:
				kim.Status.AdmissionController.State = StateResourcesCreated
			}

			if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
				xlog.WithFields(logrus.Fields{
					"namespace": deployment.Namespace,
					"name":      deployment.Name,
					"kind":      deployment.Kind,
				}).Warn("The deployment is not ready yet")
				rescheduleAfter = true
				continue
			} else {
				switch deployment.Name {
				case KimupControllerName:
					kim.Status.Controller.State = StateReady
				case KimupAdmissionControllerName:
					kim.Status.AdmissionController.State = StateReady
				}
				xlog.WithFields(logrus.Fields{
					"namespace": deployment.Namespace,
					"name":      deployment.Name,
					"kind":      deployment.Kind,
				}).Info("The deployment is ready")
			}

		case "DaemonSet":
			var daemonset appsv1.DaemonSet
			if err := r.Get(ctx, client.ObjectKeyFromObject(resource.obj), &daemonset); err != nil {
				xlog.
					WithError(err).
					WithFields(logrus.Fields{
						"namespace": daemonset.Namespace,
						"name":      daemonset.Name,
						"kind":      daemonset.Kind,
					}).Error("could not get the daemonset")
				rescheduleAfter = true
				continue
			}

			switch daemonset.Name {
			case KimupControllerName:
				kim.Status.Controller.State = StateResourcesCreated
			case KimupAdmissionControllerName:
				kim.Status.AdmissionController.State = StateResourcesCreated
			}

			if daemonset.Status.DesiredNumberScheduled != daemonset.Status.NumberReady {
				xlog.WithFields(logrus.Fields{
					"namespace": daemonset.Namespace,
					"name":      daemonset.Name,
					"kind":      daemonset.Kind,
				}).Warn("The daemonset is not ready yet")
				rescheduleAfter = true
				continue
			} else {
				switch daemonset.Name {
				case KimupControllerName:
					kim.Status.Controller.State = StateReady
				case KimupAdmissionControllerName:
					kim.Status.AdmissionController.State = StateReady
				}
				xlog.WithFields(logrus.Fields{
					"namespace": daemonset.Namespace,
					"name":      daemonset.Name,
					"kind":      daemonset.Kind,
				}).Info("The daemonset is ready")
			}

		case "Service":
			var service corev1.Service
			if err := r.Get(ctx, client.ObjectKeyFromObject(resource.obj), &service); err != nil {
				xlog.
					WithError(err).
					WithFields(logrus.Fields{
						"namespace": service.Namespace,
						"name":      service.Name,
						"kind":      service.Kind,
					}).Error("could not get the service")
				rescheduleAfter = true
				continue
			}

			if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
				if len(service.Status.LoadBalancer.Ingress) == 0 {
					xlog.WithFields(logrus.Fields{
						"namespace": service.Namespace,
						"name":      service.Name,
						"kind":      service.Kind,
					}).Warn("The service is not ready yet")
					rescheduleAfter = true
					continue
				}
			}

		default:
			xlog.WithFields(logrus.Fields{
				"kind": resource.kind,
			}).Warn("Unknown resource type")
		}
	}

	// Update status
	if err := r.Status().Update(ctx, &kim); err != nil {
		xlog.WithError(err).Error("could not update the status of the kimup object")
	}

	if rescheduleAfter {
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KimupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Listen only to spec changes
		For(&v1alpha1.Kimup{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
