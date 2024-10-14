package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

func GetKimupControllerResources(ctx context.Context, ki v1alpha1.Kimup) []Object {
	// log := log.FromContext(ctx)

	var resources []Object

	var (
		name  = KimupControllerName
		image = ki.Spec.Controller.Image
	)

	if image == "" {
		image = fmt.Sprintf("%s:%s", KimupControllerImage, Version)
	}

	// Create a deployment
	deployment := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
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
								return buildContainerPorts(ki.Spec.Controller.KimupExtraSpec)
							}(),
							Args: func() []string {
								return buildKimupArgs(ki.Spec.Controller.KimupExtraSpec)
							}(),
							ReadinessProbe:  buildReadinessProbe(ki.Spec.Controller.KimupExtraSpec),
							LivenessProbe:   buildLivenessProbe(ki.Spec.Controller.KimupExtraSpec),
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

	resources = append(resources, Object{kind: deployment.TypeMeta.Kind, obj: &deployment})

	if ki.Spec.Controller.Healthz.Enabled || ki.Spec.Controller.Metrics.Enabled {
		service := corev1.Service{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
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

		resources = append(resources, Object{kind: service.TypeMeta.Kind, obj: &service})
	}

	return resources
}
