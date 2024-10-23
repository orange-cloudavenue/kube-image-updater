package main

import (
	"context"
	"crypto/tls"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/orange-cloudavenue/kube-image-updater/internal/httpserver"
	client "github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
)

var (
	insideCluster bool = true // running inside k8s cluster

	webhookNamespace   string = "nip.io"
	webhookServiceName string = "192-168-1-30"
	webhookConfigName  string = "mutating-webhook-configuration"
	webhookPathMutate  string = "/mutate"
	webhookPort        string = ":8443"
	webhookBase               = webhookServiceName + "." + webhookNamespace

	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	kubeClient          client.Interface
	manifestWebhookPath string = "./examples/mutatingWebhookConfiguration.yaml"
)

func init() {
	// Init Metrics
	metrics.AdmissionController()

	// webhook server running namespace (default to "default")
	if os.Getenv("POD_NAMESPACE") != "" {
		webhookNamespace = os.Getenv("POD_NAMESPACE")
	}
	// init flags
	flag.StringVar(&webhookPort, "webhook-port", webhookPort, "Webhook server port.ex: :8443")
	flag.StringVar(&webhookNamespace, "namespace", webhookNamespace, "Kimup Webhook Mutating namespace.")
	flag.StringVar(&webhookServiceName, "service-name", webhookServiceName, "Kimup Webhook Mutating service name.")
	flag.BoolVar(&insideCluster, "inside-cluster", true, "True if running inside k8s cluster.")
	flag.Parse()
}

// Start http server for webhook
func main() {
	var err error

	// -- Context -- //
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// -- OS signal handling -- //
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	// kubernetes golang library provide flag "kubeconfig" to specify the path to the kubeconfig file
	kubeClient, err = client.New(flag.Lookup("kubeconfig").Value.String(), client.ComponentAdmissionController)
	if err != nil {
		log.WithError(err).Panicf("Error creating kubeclient")
	}

	// * Webhook server
	// generate cert for webhook
	pair, caPEM, err := generateTLS()
	if err != nil {
		log.WithError(err).Fatal("Failed to generate TLS")
	}
	tlsC := &tls.Config{
		Certificates: []tls.Certificate{pair},
		MinVersion:   tls.VersionTLS12,
		// InsecureSkipVerify: true, //nolint:gosec
	}

	// create or update the mutatingwebhookconfiguration
	err = createOrUpdateMutatingWebhookConfiguration(caPEM, webhookServiceName, webhookNamespace, kubeClient)
	if err != nil {
		log.WithError(err).Error("Failed to create or update the mutating webhook configuration")
		signalChan <- os.Interrupt
	}

	// * Config the webhook server
	a, waitHTTP := httpserver.Init(ctx, httpserver.WithCustomHandlerForHealth(
		func() (bool, error) {
			_, err := net.DialTimeout("tcp", webhookPort, 5*time.Second)
			if err != nil {
				return false, err
			}
			return true, nil
		}))

	s, err := a.Add("webhook", httpserver.WithTLS(tlsC), httpserver.WithAddr(webhookPort))
	if err != nil {
		log.
			WithError(err).
			WithFields(logrus.Fields{
				"address": webhookPort,
			}).Fatal("Failed to create the server")
	}
	s.Config.Post(webhookPathMutate, ServeHandler)
	if err := a.Run(); err != nil {
		log.WithError(err).Fatal("Failed to start HTTP servers")
	}

	// !-- OS signal handling --! //
	<-signalChan
	// cancel the context
	cancel()
	waitHTTP()
}
