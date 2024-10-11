package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
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
func (r *KimupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var kim v1alpha1.Kimup
	if err := r.Get(ctx, req.NamespacedName, &kim); err != nil {
		log.Info(fmt.Sprintf("cloud not get the Kimup object: %s", req.NamespacedName))
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling kimup objects")
	if kim.Spec.Controller != nil {
		log.Info("Reconciling kimup controller")

		switch {
		// if phase is empty, create the resources
		case kim.Status.Controller.Phase == "":
			log.Info("Creating the resources for the controller")
			resources, err := GetKimupControllerResources(ctx, kim)
			if err != nil {
				log.Error(err, "could not get the resources for the controller")
				return ctrl.Result{}, err
			}

			for _, resource := range resources {
				if err := r.Create(ctx, resource); err != nil {
					log.Error(err, "could not create the resource")
					return ctrl.Result{RequeueAfter: 10 * time.Second}, err
				}
			}

			// Update status
			kim.Status.Controller.Phase = PhaseResourcesCreated
			log.Info("Created resources for the controller")
			if err := r.Status().Update(ctx, &kim); err != nil {
				log.Error(err, "could not update the status of the kimup controller")
				return ctrl.Result{}, err
			}

			r.Recorder.Event(&kim, corev1.EventTypeNormal, "Resources", "Created resources for the controller")
			return ctrl.Result{}, nil
		default:
			var deployment appsv1.Deployment
			if err := r.Get(ctx, client.ObjectKey{Namespace: kim.Namespace, Name: KimupControllerName}, &deployment); err != nil {
				log.Error(err, "could not get the deployment for the controller")
				return ctrl.Result{}, err
			}

			if deployment.Status.Replicas != deployment.Status.ReadyReplicas {
				log.Info("The controller is not ready yet")
				return ctrl.Result{
					Requeue: true,
				}, nil
			}

			resources, err := GetKimupControllerResources(ctx, kim)
			if err != nil {
				log.Error(err, "could not get the resources for the controller")
				return ctrl.Result{}, err
			}

			for _, resource := range resources {
				if err := r.Update(ctx, resource); err != nil {
					log.Error(err, "could not update the resource %s/%s/%s", resource.GetObjectKind(), resource.GetNamespace(), resource.GetName())
					return ctrl.Result{}, err
				}
			}

			log.Info("Updated resources for the controller")
			r.Recorder.Event(&kim, corev1.EventTypeNormal, "Resources", "Updated resources for the controller")
			return ctrl.Result{}, nil
		}
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

func GetKimupControllerResources(ctx context.Context, ki v1alpha1.Kimup) ([]client.Object, error) {
	// log := log.FromContext(ctx)

	var resources []client.Object

	var (
		name  = KimupControllerName
		image = ki.Spec.Controller.Image
	)

	if image == "" {
		image = fmt.Sprintf("%s:%s", KimupControllerImage, Version)
	}

	// Create a deployment
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ki.Namespace,
			// Useful for automatically deleting the resources when the kimup object is deleted
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ki.APIVersion,
					Kind:       ki.Kind,
					Name:       ki.Name,
					UID:        ki.UID,
				},
			},
			Labels: map[string]string{
				KubernetesAppComponentLabelKey: KimupControllerName,
				KubernetesAppInstanceNameLabel: name,
				KubernetesAppNameLabelKey:      KimupControllerName,
				KubernetesAppVersionLabelKey:   Version,
				KubernetesPartOfLabelKey:       KimupControllerName,
				KubernetesManagedByLabelKey:    KimupOperatorName,
				"app":                          name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.ToPTR(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":                     name,
					KubernetesPartOfLabelKey:  KimupControllerName,
					KubernetesAppNameLabelKey: KimupControllerName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: func() map[string]string {
						labels := map[string]string{
							"app":                     name,
							KubernetesPartOfLabelKey:  KimupControllerName,
							KubernetesAppNameLabelKey: KimupControllerName,
						}
						for k, v := range ki.Spec.Controller.Labels {
							labels[k] = v
						}
						return labels
					}(),
					Annotations: ki.Spec.Controller.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kimup",
							Image: image,
							Ports: func() []corev1.ContainerPort {
								ports := []corev1.ContainerPort{}

								if ki.Spec.Controller.Metrics.Enabled {
									// set the metrics port
									metricsPort := ki.Spec.Controller.Metrics.Port
									if metricsPort != 0 {
										metricsPort = models.MetricsDefaultPort
									}

									ports = append(ports, corev1.ContainerPort{
										Name:          models.MetricsFlagName,
										ContainerPort: metricsPort,
									})
								}

								if ki.Spec.Controller.Healthz.Enabled {
									// set the healthz port
									healthzPort := ki.Spec.Controller.Healthz.Port
									if healthzPort != 0 {
										healthzPort = models.HealthzDefaultPort
									}

									ports = append(ports, corev1.ContainerPort{
										Name:          models.HealthzFlagName,
										ContainerPort: healthzPort,
									})
								}

								return ports
							}(),
							Args: func() []string {
								a := []string{}
								if ki.Spec.Controller.Healthz.Enabled {
									// enable healthz
									a = append(a, fmt.Sprintf("--%s", models.HealthzFlagName))

									// set the healthz port
									healthzPort := ki.Spec.Controller.Healthz.Port
									if healthzPort != 0 {
										healthzPort = models.HealthzDefaultPort
									}
									a = append(a, fmt.Sprintf("--%s=%d", models.HealthzPortFlagName, healthzPort))

									// set the healthz path
									healthzPath := ki.Spec.Controller.Healthz.Path
									if healthzPath == "" {
										healthzPath = models.HealthzDefaultPath
									}
									a = append(a, fmt.Sprintf("--%s=%s", models.HealthzPathFlagName, healthzPath))
								}

								if ki.Spec.Controller.Metrics.Enabled {
									// enable metrics
									a = append(a, fmt.Sprintf("--%s", models.MetricsFlagName))

									// set the metrics port
									metricsPort := ki.Spec.Controller.Metrics.Port
									if metricsPort != 0 {
										metricsPort = models.MetricsDefaultPort
									}

									a = append(a, fmt.Sprintf("--%s=%d", models.MetricsPortFlagName, metricsPort))

									// set the metrics path
									metricsPath := ki.Spec.Controller.Metrics.Path
									if metricsPath == "" {
										metricsPath = models.MetricsDefaultPath
									}

									a = append(a, fmt.Sprintf("--%s=%s", models.MetricsPathFlagName, metricsPath))
								}

								a = append(a, fmt.Sprintf("--%s=%s", models.LogLevelFlagName, ki.Spec.Controller.LogLevel))

								return a
							}(),
							ReadinessProbe: func() *corev1.Probe {
								if !ki.Spec.Controller.Healthz.Enabled {
									return nil
								}
								healthzPath := ki.Spec.Controller.Healthz.Path
								if healthzPath == "" {
									healthzPath = models.HealthzDefaultPath
								}

								healthzPort := ki.Spec.Controller.Healthz.Port
								if healthzPort != 0 {
									healthzPort = models.HealthzDefaultPort
								}

								return &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: healthzPath,
											Port: intstr.FromInt32(healthzPort),
										},
									},
									FailureThreshold:    3,
									InitialDelaySeconds: 10,
									PeriodSeconds:       10,
									SuccessThreshold:    1,
									TimeoutSeconds:      2,
								}
							}(),
							LivenessProbe: func() *corev1.Probe {
								if !ki.Spec.Controller.Healthz.Enabled {
									return nil
								}

								healthzPath := ki.Spec.Controller.Healthz.Path
								if healthzPath == "" {
									healthzPath = models.HealthzDefaultPath
								}

								healthzPort := ki.Spec.Controller.Healthz.Port
								if healthzPort != 0 {
									healthzPort = models.HealthzDefaultPort
								}

								return &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: healthzPath,
											Port: intstr.FromInt32(healthzPort),
										},
									},
									FailureThreshold:    3,
									InitialDelaySeconds: 10,
									PeriodSeconds:       10,
									SuccessThreshold:    1,
									TimeoutSeconds:      2,
								}
							}(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Resources: func() corev1.ResourceRequirements {
								if ki.Spec.Controller.Resources == nil {
									return corev1.ResourceRequirements{}
								}
								return *ki.Spec.Controller.Resources
							}(),
						},
					},
					Affinity:                  ki.Spec.Controller.Affinity,
					NodeSelector:              ki.Spec.Controller.NodeSelector,
					Tolerations:               ki.Spec.Controller.Tolerations,
					TopologySpreadConstraints: ki.Spec.Controller.TopologySpreadConstraints,
					ServiceAccountName:        ki.Spec.Controller.ServiceAccountName,
					PriorityClassName:         ki.Spec.Controller.PriorityClassName,
				},
			},
		},
	}

	resources = append(resources, &deployment)

	// TODO Check if healthz or metrics are enabled

	service := corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-healthz-metrics", name),
			Namespace: ki.Namespace,
			// Useful for automatically deleting the resources when the kimup object is deleted
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: ki.APIVersion,
					Kind:       ki.Kind,
					Name:       ki.Name,
					UID:        ki.UID,
				},
			},
			Labels: map[string]string{
				KubernetesAppComponentLabelKey: KimupControllerName,
				KubernetesAppInstanceNameLabel: name,
				KubernetesAppNameLabelKey:      KimupControllerName,
				KubernetesAppVersionLabelKey:   Version,
				KubernetesPartOfLabelKey:       KimupControllerName,
				KubernetesManagedByLabelKey:    KimupOperatorName,
				"app":                          name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app":                     name,
				KubernetesAppNameLabelKey: KimupControllerName,
			},
			Ports: func() []corev1.ServicePort {
				svcs := []corev1.ServicePort{}

				if ki.Spec.Controller.Metrics.Enabled {
					// set the metrics port
					metricsPort := ki.Spec.Controller.Metrics.Port
					if metricsPort == 0 {
						metricsPort = models.MetricsDefaultPort
					}

					svcs = append(svcs, corev1.ServicePort{
						Name:       models.MetricsFlagName,
						Port:       metricsPort,
						TargetPort: intstr.FromString(models.MetricsFlagName),
					})
				}

				if ki.Spec.Controller.Healthz.Enabled {
					// set the healthz port
					healthzPort := ki.Spec.Controller.Healthz.Port
					if healthzPort == 0 {
						healthzPort = models.HealthzDefaultPort
					}

					svcs = append(svcs, corev1.ServicePort{
						Name:       models.HealthzFlagName,
						Port:       healthzPort,
						TargetPort: intstr.FromString(models.HealthzFlagName),
					})
				}

				return svcs
			}(),
		},
	}

	resources = append(resources, &service)

	return resources, nil
}
