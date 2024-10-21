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

package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bombsimon/logrusr/v4"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	kimupv1alpha1 "github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/controller"
	"github.com/orange-cloudavenue/kube-image-updater/internal/httpserver"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	scheme = runtime.NewScheme()
	c      = make(chan os.Signal, 1)
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(kimupv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	var enableLeaderElection bool
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(logrusr.New(log.GetLogger()))

	webhook := webhook.NewServer(webhook.Options{})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			// TODO pass metrics to common metrics server
			BindAddress: "0", // disable metrics service
		},
		// TODO reenable and use comman healthz server to check the health of the application
		HealthProbeBindAddress: "0", // disable health probe service
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "71be4586.cloudavenue.io",
		WebhookServer:          webhook,
	})
	if err != nil {
		log.WithError(err).Error("unable to start manager")
		c <- syscall.SIGINT
	}

	kubeAPIClient, err := kubeclient.NewFromRestConfig(ctrl.GetConfigOrDie(), kubeclient.ComponentOperator)
	if err != nil {
		log.WithError(err).Error("unable to create kubeclient")
		c <- syscall.SIGINT
	}

	// ! Mutator

	if err := (&controller.ImageTagMutator{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		KubeAPIClient: kubeAPIClient,
	}).SetupWebhookWithManager(mgr); err != nil {
		log.WithError(err).Error("unable to create webhook", "webhook", "ImageTagMutator")
		c <- syscall.SIGINT
	}

	// ! Reconcilers

	if err = (&controller.ImageReconciler{
		Client:        mgr.GetClient(),
		KubeAPIClient: kubeAPIClient,
		Scheme:        mgr.GetScheme(),
		Recorder:      mgr.GetEventRecorderFor("kimup-operator"),
	}).SetupWithManager(mgr); err != nil {
		log.WithError(err).Error("unable to create controller", "controller", "Image")
		c <- syscall.SIGINT
	}

	if err = (&controller.KimupReconciler{
		Client:        mgr.GetClient(),
		KubeAPIClient: kubeAPIClient,
		Scheme:        mgr.GetScheme(),
		Recorder:      mgr.GetEventRecorderFor("kimup-operator"),
	}).SetupWithManager(mgr); err != nil {
		log.WithError(err).Error(err, "unable to create controller", "controller", "Kimup")
		c <- syscall.SIGINT
	}

	// +kubebuilder:scaffold:builder

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// * Config the metrics and healthz server
	a, waitHTTP := httpserver.Init(ctx, httpserver.WithCustomHandlerForHealth(
		func() (bool, error) {
			// TODO improve
			_, err := net.DialTimeout("tcp", models.HealthzDefaultAddr, 5*time.Second)
			if err != nil {
				return false, err
			}
			return true, nil
		}))

	if err := a.Run(); err != nil {
		log.WithError(err).Error("Failed to start HTTP servers")
		// send signal to stop the program
		c <- syscall.SIGINT
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.WithError(err).Error(err, "unable to set up health check")
		c <- syscall.SIGINT
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.WithError(err).Error(err, "unable to set up ready check")
		c <- syscall.SIGINT
	}

	log.Info("Starting operator")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.WithError(err).Error(err, "problem running manager")
		c <- syscall.SIGINT
	}

	<-c
	cancel()
	waitHTTP()
}
