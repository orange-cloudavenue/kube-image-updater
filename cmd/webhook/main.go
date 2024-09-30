package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

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

	webhookMetricsPath string = "/metrics"
	webhookMetricsPort string = ":9080"

	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	kubeClient          *client.Client
	manifestWebhookPath string = "./config/manifests/mutatingWebhookConfiguration.yaml"

	// Prometheus metrics
	promHTTPRequestsTotal prometheus.Counter   = metrics.NewCounter("http_requests_total", "The total number of handled HTTP requests.")
	promHTTPErrorsTotal   prometheus.Counter   = metrics.NewCounter("http_errors_total", "The total number of handled HTTP errors.")
	promHTTPDuration      prometheus.Histogram = metrics.NewHistogram("http_response_time_seconds", "The duration in seconds of HTTP requests.")
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
}

// Start http server for webhook
func main() {
	var err error
	flag.StringVar(&webhookPort, "webhook-port", webhookPort, "Webhook server port.ex: :8443")
	flag.StringVar(&webhookNamespace, "namespace", webhookNamespace, "Kimup Webhook Mutating namespace.")
	flag.StringVar(&webhookServiceName, "service-name", webhookServiceName, "Kimup Webhook Mutating service name.")

	flag.BoolVar(&insideCluster, "inside-cluster", true, "True if running inside k8s cluster.")

	flag.Parse()

	// homedir for kubeconfig
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	kubeClient, err = client.New(homedir + "/.kube/config")
	if err != nil {
		panic(err)
	}

	// generate cert for webhook
	pair, caPEM := generateTLS()

	mux := http.NewServeMux()
	mux.HandleFunc(webhookPathMutate, serveHandler)

	// create or update the mutatingwebhookconfiguration
	err = createOrUpdateMutatingWebhookConfiguration(caPEM, webhookServiceName, webhookNamespace, kubeClient)
	if err != nil {
		errorLogger.Fatalf("Failed to create or update the mutating webhook configuration: %v", err)
	}

	// define http server and server handler
	s := &http.Server{
		Addr:        webhookPort,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{pair},
			MinVersion:   tls.VersionTLS12,
			// InsecureSkipVerify: true, //nolint:gosec
		},
	}

	// start the HTTP server
	go func() {
		infoLogger.Printf("Starting webhook server on %s from insideCluster=%v", s.Addr, insideCluster)
		if err = s.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("Failed to start webhook server: %v", err)
		}
	}()

	if err := metrics.ServeProm(webhookMetricsPort, webhookMetricsPath); err != nil {
		errorLogger.Fatalf("Failed to start metrics server: %v", err)
	}

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	infoLogger.Printf("Got OS shutdown signal, shutting down webhook server gracefully...")
	s.Shutdown(context.Background()) //nolint:errcheck
}
