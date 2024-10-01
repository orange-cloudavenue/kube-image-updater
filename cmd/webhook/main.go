package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net"
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
	webhookHost        string = "192.168.1.48"

	webhookMetricsPath string = "/metrics"
	webhookMetricsPort string = ":9080"

	webhookHealthPort string = ":9081"
	webhookHealthPath string = "/healthz"

	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	kubeClient          *client.Client
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
}

// Start http server for webhook
func main() {
	var err error
	flag.StringVar(&webhookPort, "webhook-port", webhookPort, "Webhook server port.ex: :8443")
	flag.StringVar(&webhookNamespace, "namespace", webhookNamespace, "Kimup Webhook Mutating namespace.")
	flag.StringVar(&webhookServiceName, "service-name", webhookServiceName, "Kimup Webhook Mutating service name.")
	flag.StringVar(&webhookHost, "webhook-host", webhookHost, "Webhook server host.")

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

	// !-- Webhook server --! //
	// generate cert for webhook
	pair, caPEM := generateTLS()

	// create or update the mutatingwebhookconfiguration
	err = createOrUpdateMutatingWebhookConfiguration(caPEM, webhookServiceName, webhookNamespace, kubeClient)
	if err != nil {
		errorLogger.Fatalf("Failed to create or update the mutating webhook configuration: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(webhookPathMutate, serveHandler)

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

	// !-- Prometheus metrics server --! //
	// start the metrics server
	if err := metrics.ServeProm(webhookMetricsPort, webhookMetricsPath); err != nil {
		errorLogger.Fatalf("Failed to start metrics server: %v", err)
	}

	// !-- Health check server --! //
	// define health check
	muxHealth := http.NewServeMux()
	muxHealth.HandleFunc(webhookHealthPath, func(w http.ResponseWriter, r *http.Request) {
		_, err := net.DialTimeout("tcp", webhookPort, 1*time.Second)
		if err != nil {
			log.Println("Webhook port is unreachable, error: ", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			errorLogger.Printf("Failed to write health check response: %v", err)
		}
	})

	// create health check server
	sc := &http.Server{
		Addr:        webhookHealthPort,
		Handler:     muxHealth,
		ReadTimeout: 10 * time.Second,
	}

	go func() {
		infoLogger.Printf("Starting health check server on %s", sc.Addr)
		if err = sc.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start health check server: %v", err)
		}
	}()

	// !-- OS signal handling --! //
	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	infoLogger.Printf("Got OS shutdown signal, shutting down webhook server gracefully...")
	s.Shutdown(context.Background()) //nolint:errcheck
}
