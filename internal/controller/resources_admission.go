package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

func GetKimupAdmissionResources(ctx context.Context, ki v1alpha1.Kimup) []Object {
	// log := log.FromContext(ctx)

	var resources []Object

	var (
		name  = KimupAdmissionControllerName
		image = ki.Spec.AdmissionController.Image
	)

	if image == "" {
		image = fmt.Sprintf("%s:%s", KimupAdmissionControllerImage, Version)
	}

	switch ki.Spec.AdmissionController.DeploymentType {
	case "Deployment":
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
					KubernetesAppComponentLabelKey: name,
					KubernetesAppInstanceNameLabel: name,
					KubernetesAppNameLabelKey:      name,
					KubernetesAppVersionLabelKey:   Version,
					KubernetesPartOfLabelKey:       name,
					KubernetesManagedByLabelKey:    KimupOperatorName,
					"app":                          name,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &ki.Spec.AdmissionController.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":                     name,
						KubernetesPartOfLabelKey:  name,
						KubernetesAppNameLabelKey: name,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: func() map[string]string {
							labels := map[string]string{
								"app":                     name,
								KubernetesPartOfLabelKey:  name,
								KubernetesAppNameLabelKey: name,
							}
							for k, v := range ki.Spec.AdmissionController.Labels {
								labels[k] = v
							}
							return labels
						}(),
						Annotations: ki.Spec.AdmissionController.Annotations,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "kimup-admission-controller",
								Image: image,
								Ports: func() []corev1.ContainerPort {
									return buildContainerPorts(ki.Spec.AdmissionController.KimupExtraSpec)
								}(),
								Args: func() []string {
									return buildKimupArgs(ki.Spec.AdmissionController.KimupExtraSpec)
								}(),
								ReadinessProbe:  buildReadinessProbe(ki.Spec.AdmissionController.KimupExtraSpec),
								LivenessProbe:   buildLivenessProbe(ki.Spec.AdmissionController.KimupExtraSpec),
								ImagePullPolicy: corev1.PullIfNotPresent,
								Resources: func() corev1.ResourceRequirements {
									if ki.Spec.AdmissionController.Resources == nil {
										return corev1.ResourceRequirements{}
									}
									return *ki.Spec.AdmissionController.Resources
								}(),
							},
						},
						Affinity:                  ki.Spec.AdmissionController.Affinity,
						NodeSelector:              ki.Spec.AdmissionController.NodeSelector,
						Tolerations:               ki.Spec.AdmissionController.Tolerations,
						TopologySpreadConstraints: ki.Spec.AdmissionController.TopologySpreadConstraints,
						ServiceAccountName:        ki.Spec.AdmissionController.ServiceAccountName,
						PriorityClassName:         ki.Spec.AdmissionController.PriorityClassName,
					},
				},
			},
		}

		resources = append(resources, Object{kind: deployment.TypeMeta.Kind, obj: &deployment})
	case "DaemonSet":
		// Create a daemonset

		daemonset := appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "DaemonSet",
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
					KubernetesAppComponentLabelKey: name,
					KubernetesAppInstanceNameLabel: name,
					KubernetesAppNameLabelKey:      name,
					KubernetesAppVersionLabelKey:   Version,
					KubernetesPartOfLabelKey:       name,
					KubernetesManagedByLabelKey:    KimupOperatorName,
					"app":                          name,
				},
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app":                     name,
						KubernetesPartOfLabelKey:  name,
						KubernetesAppNameLabelKey: name,
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: func() map[string]string {
							labels := map[string]string{
								"app":                     name,
								KubernetesPartOfLabelKey:  name,
								KubernetesAppNameLabelKey: name,
							}
							for k, v := range ki.Spec.AdmissionController.Labels {
								labels[k] = v
							}
							return labels
						}(),
						Annotations: ki.Spec.AdmissionController.Annotations,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "kimup",
								Image: image,
								Ports: func() []corev1.ContainerPort {
									return buildContainerPorts(ki.Spec.AdmissionController.KimupExtraSpec)
								}(),
								Args: func() []string {
									return buildKimupArgs(ki.Spec.AdmissionController.KimupExtraSpec)
								}(),
								ReadinessProbe:  buildReadinessProbe(ki.Spec.AdmissionController.KimupExtraSpec),
								LivenessProbe:   buildLivenessProbe(ki.Spec.AdmissionController.KimupExtraSpec),
								ImagePullPolicy: corev1.PullIfNotPresent,
								Resources: func() corev1.ResourceRequirements {
									if ki.Spec.AdmissionController.Resources == nil {
										return corev1.ResourceRequirements{}
									}
									return *ki.Spec.AdmissionController.Resources
								}(),
							},
						},
						Affinity:                  ki.Spec.AdmissionController.Affinity,
						NodeSelector:              ki.Spec.AdmissionController.NodeSelector,
						Tolerations:               ki.Spec.AdmissionController.Tolerations,
						TopologySpreadConstraints: ki.Spec.AdmissionController.TopologySpreadConstraints,
						ServiceAccountName:        ki.Spec.AdmissionController.ServiceAccountName,
						PriorityClassName:         ki.Spec.AdmissionController.PriorityClassName,
					},
				},
			},
		}
		resources = append(resources, Object{kind: daemonset.TypeMeta.Kind, obj: &daemonset})
	} // end switch

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
					return buildServicePorts(ki.Spec.AdmissionController.KimupExtraSpec)
				}(),
			},
		}

		resources = append(resources, Object{kind: service.TypeMeta.Kind, obj: &service})
	}
	return resources
}
