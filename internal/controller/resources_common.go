package controller

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

type Object struct {
	kind string
	obj  client.Object
}

func buildKimupArgs(extra v1alpha1.KimupExtraSpec) (args []string) {
	args = []string{}
	if extra.Healthz.Enabled {
		// enable healthz
		args = append(args, fmt.Sprintf("--%s", models.HealthzFlagName))

		// set the healthz port
		healthzPort := extra.Healthz.Port
		if healthzPort != 0 {
			healthzPort = models.HealthzDefaultPort
		}
		args = append(args, fmt.Sprintf("--%s=%d", models.HealthzPortFlagName, healthzPort))

		// set the healthz path
		healthzPath := extra.Healthz.Path
		if healthzPath == "" {
			healthzPath = models.HealthzDefaultPath
		}
		args = append(args, fmt.Sprintf("--%s=%s", models.HealthzPathFlagName, healthzPath))
	}

	if extra.Metrics.Enabled {
		// enable metrics
		args = append(args, fmt.Sprintf("--%s", models.MetricsFlagName))

		// set the metrics port
		metricsPort := extra.Metrics.Port
		if metricsPort != 0 {
			metricsPort = models.MetricsDefaultPort
		}

		args = append(args, fmt.Sprintf("--%s=%d", models.MetricsPortFlagName, metricsPort))

		// set the metrics path
		metricsPath := extra.Metrics.Path
		if metricsPath == "" {
			metricsPath = models.MetricsDefaultPath
		}

		args = append(args, fmt.Sprintf("--%s=%s", models.MetricsPathFlagName, metricsPath))
	}

	// TODO
	args = append(args, fmt.Sprintf("--%s=%s", models.LogLevelFlagName, extra.LogLevel))

	return args
}

func buildReadinessProbe(extra v1alpha1.KimupExtraSpec) *corev1.Probe {
	if !extra.Healthz.Enabled {
		return nil
	}
	healthzPath := extra.Healthz.Path
	if healthzPath == "" {
		healthzPath = models.HealthzDefaultPath
	}

	healthzPort := extra.Healthz.Port
	if healthzPort == 0 {
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
}

func buildLivenessProbe(extra v1alpha1.KimupExtraSpec) *corev1.Probe {
	return buildReadinessProbe(extra)
}

func buildContainerPorts(extra v1alpha1.KimupExtraSpec) (ports []corev1.ContainerPort) {
	ports = []corev1.ContainerPort{}

	if extra.Metrics.Enabled {
		// set the metrics port
		metricsPort := extra.Metrics.Port
		if metricsPort == 0 {
			metricsPort = models.MetricsDefaultPort
		}

		ports = append(ports, corev1.ContainerPort{
			Name:          models.MetricsFlagName,
			ContainerPort: metricsPort,
		})
	}

	if extra.Healthz.Enabled {
		// set the healthz port
		healthzPort := extra.Healthz.Port
		if healthzPort == 0 {
			healthzPort = models.HealthzDefaultPort
		}

		ports = append(ports, corev1.ContainerPort{
			Name:          models.HealthzFlagName,
			ContainerPort: healthzPort,
		})
	}

	return ports
}

func buildServicePorts(extra v1alpha1.KimupExtraSpec) (ports []corev1.ServicePort) {
	ports = []corev1.ServicePort{}

	if extra.Metrics.Enabled {
		// set the metrics port
		metricsPort := extra.Metrics.Port
		if metricsPort == 0 {
			metricsPort = models.MetricsDefaultPort
		}

		ports = append(ports, corev1.ServicePort{
			Name:       models.MetricsFlagName,
			Port:       metricsPort,
			TargetPort: intstr.FromString(models.MetricsFlagName),
		})
	}

	if extra.Healthz.Enabled {
		// set the healthz port
		healthzPort := extra.Healthz.Port
		if healthzPort == 0 {
			healthzPort = models.HealthzDefaultPort
		}

		ports = append(ports, corev1.ServicePort{
			Name:       models.HealthzFlagName,
			Port:       healthzPort,
			TargetPort: intstr.FromString(models.HealthzFlagName),
		})
	}

	return ports
}
