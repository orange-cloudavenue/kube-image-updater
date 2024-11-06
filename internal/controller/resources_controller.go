package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

func GetKimupControllerResources(ctx context.Context, ki v1alpha1.Kimup) []Object {
	// log := log.FromContext(ctx)

	var resources []Object

	var (
		name  = KimupControllerName
		image = ki.Spec.Image
	)

	if image == "" {
		image = fmt.Sprintf("%s:%s", KimupControllerImage, models.Version)
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
				KubernetesAppVersionLabelKey:   models.Version,
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
						for k, v := range ki.Spec.Labels {
							labels[k] = v
						}
						return labels
					}(),
					Annotations: ki.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kimup",
							Image: image,
							Ports: func() []corev1.ContainerPort {
								return buildContainerPorts(ki.Spec.KimupExtraSpec)
							}(),
							Args: func() []string {
								return buildKimupArgs(ki.Spec.KimupExtraSpec)
							}(),
							ReadinessProbe:  buildReadinessProbe(ki.Spec.KimupExtraSpec),
							LivenessProbe:   buildLivenessProbe(ki.Spec.KimupExtraSpec),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Resources: func() corev1.ResourceRequirements {
								if ki.Spec.Resources == nil {
									return corev1.ResourceRequirements{}
								}
								return *ki.Spec.Resources
							}(),
						},
					},
					Affinity:                  ki.Spec.Affinity,
					NodeSelector:              ki.Spec.NodeSelector,
					Tolerations:               ki.Spec.Tolerations,
					TopologySpreadConstraints: ki.Spec.TopologySpreadConstraints,
					ServiceAccountName:        ki.Spec.ServiceAccountName,
					PriorityClassName:         ki.Spec.PriorityClassName,
				},
			},
		},
	}

	resources = append(resources, Object{kind: deployment.TypeMeta.Kind, obj: &deployment})

	if ki.Spec.Healthz.Enabled || ki.Spec.Metrics.Enabled {
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
					KubernetesAppVersionLabelKey:   models.Version,
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
				Ports: buildServicePorts(ki.Spec.KimupExtraSpec),
			},
		}

		resources = append(resources, Object{kind: service.TypeMeta.Kind, obj: &service})
	}

	return resources
}
