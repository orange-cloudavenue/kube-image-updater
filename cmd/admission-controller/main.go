package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/orange-cloudavenue/kube-image-updater/internal/httpserver"
	client "github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
)

var (
	insideCluster bool = true // running inside k8s cluster
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	webhookNamespace   string = "example.com"
	webhookServiceName string = "your"
	webhookConfigName  string = "webhookconfig"
	webhookPathMutate  string = "/mutate"
	webhookPort        string = ":8443"
	webhookBase               = webhookServiceName + "." + webhookNamespace

	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	kubeClient          client.Interface
	manifestWebhookPath string = "./config/manifests/mutatingWebhookConfiguration.yaml"

	// Prometheus metrics
	promHTTPRequestsTotal prometheus.Counter   = metrics.NewCounter("http_requests_total", "The total number of handled HTTP requests.")
	promHTTPErrorsTotal   prometheus.Counter   = metrics.NewCounter("http_errors_total", "The total number of handled HTTP errors.")
	promHTTPDuration      prometheus.Histogram = metrics.NewHistogram("http_response_time_seconds", "The duration in seconds of HTTP requests.")
	promPatchTotal        prometheus.Counter   = metrics.NewCounter("patch_total", "The total number of requests to a patch.")
)

func init() {
	// init loggers
	debugLogger = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

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
	kubeClient, err = client.New(flag.Lookup("kubeconfig").Value.String())
	if err != nil {
		log.Panicf("Error creating kubeclient: %v", err)
	}

	// * Webhook server
	// generate cert for webhook
	pair, caPEM, err := generateTLS()
	if err != nil {
		errorLogger.Fatalf("Failed to generate TLS pair: %v", err)
	}
	tlsC := &tls.Config{
		Certificates: []tls.Certificate{pair},
		MinVersion:   tls.VersionTLS12,
		// InsecureSkipVerify: true, //nolint:gosec
	}

	// create or update the mutatingwebhookconfiguration
	err = createOrUpdateMutatingWebhookConfiguration(caPEM, webhookServiceName, webhookNamespace, kubeClient)
	if err != nil {
		errorLogger.Printf("Failed to create or update the mutating webhook configuration: %v", err)
		signalChan <- os.Interrupt
	}

	// * Config the webhook server
	a, waitHTTP := httpserver.Init(ctx, httpserver.WithCustomHandlerForHealth(
		func() (bool, error) {
			_, err := net.DialTimeout("tcp", ":4444", 5*time.Second)
			if err != nil {
				return false, err
			}
			return true, nil
		}))

	s, err := a.Add("webhook", httpserver.WithTLS(tlsC), httpserver.WithAddr(webhookPort))
	if err != nil {
		errorLogger.Fatalf("Failed to create the server: %v", err)
	}
	s.Config.Post(webhookPathMutate, ServeHandler)
	if err := a.Run(); err != nil {
		errorLogger.Fatalf("Failed to start HTTP servers: %v", err)
	}

	// !-- OS signal handling --! //
	<-signalChan
	// cancel the context
	cancel()
	waitHTTP()
}
