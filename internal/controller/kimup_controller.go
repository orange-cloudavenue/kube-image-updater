package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
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
	log := log.FromContext(ctx)

	var kim v1alpha1.Kimup
	if err := r.Get(ctx, req.NamespacedName, &kim); err != nil {
		log.Info(fmt.Sprintf("cloud not get the Kimup object: %s", req.NamespacedName))
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
			log.Error(err, "could not create the resource")
			rescheduleAfter = true
			continue
		}
	}

	for _, resource := range resourcesToUpdate {
		if err := r.Update(ctx, resource.obj); err != nil {
			log.Error(err, "could not update the resource")
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
				log.Error(err, "could not get the deployment")
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
				log.Info(fmt.Sprintf("The %s deployment in namespace %s is not ready yet", deployment.Name, deployment.Namespace))
				rescheduleAfter = true
				continue
			} else {
				switch deployment.Name {
				case KimupControllerName:
					kim.Status.Controller.State = StateReady
				case KimupAdmissionControllerName:
					kim.Status.AdmissionController.State = StateReady
				}
				log.Info(fmt.Sprintf("The %s deployment in namespace %s is ready", deployment.Name, deployment.Namespace))
			}

		case "DaemonSet":
			var daemonset appsv1.DaemonSet
			if err := r.Get(ctx, client.ObjectKeyFromObject(resource.obj), &daemonset); err != nil {
				log.Error(err, "could not get the daemonset")
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
				log.Info("The daemonset is not ready yet")
				rescheduleAfter = true
				continue
			} else {
				switch daemonset.Name {
				case KimupControllerName:
					kim.Status.Controller.State = StateReady
				case KimupAdmissionControllerName:
					kim.Status.AdmissionController.State = StateReady
				}
				log.Info(fmt.Sprintf("The %s daemonset in namespace %s is ready", daemonset.Name, daemonset.Namespace))
			}

		case "Service":
			var service corev1.Service
			if err := r.Get(ctx, client.ObjectKeyFromObject(resource.obj), &service); err != nil {
				log.Error(err, "could not get the service")
				rescheduleAfter = true
				continue
			}

			if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
				if len(service.Status.LoadBalancer.Ingress) == 0 {
					log.Info("The service is not ready yet")
					rescheduleAfter = true
					continue
				}
			}

		default:
			log.Info(fmt.Sprintf("Unknown resource type: %s", resource.kind))
		}
	}

	// Update status
	if err := r.Status().Update(ctx, &kim); err != nil {
		log.Error(err, "could not update the status of the kimup object")
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
