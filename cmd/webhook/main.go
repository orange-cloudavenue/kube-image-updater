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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	client "github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
)

const (
	admissionWebhookAnnotationBase = "kimup-webhook-mutating.default.svc" // service name of the webhook
)

var (
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	webhookNamespace   string = "default"
	webhookServiceName string = "kimup-webhook-mutating"
	webhookConfigName  string = "webhookconfig"
	webhookPathMutate  string = "/mutate"
	webhookPort        string = ":8443"

	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

type (
	patchOperation struct {
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value,omitempty"`
	}
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
	flag.StringVar(&webhookServiceName, "service-name", webhookServiceName, "Kimup Webhook Mutating service name.")

	flag.Parse()

	// homedir for kubeconfig
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	k, err := client.New(homedir + "/.kube/config")
	if err != nil {
		panic(err)
	}

	// generate cert for webhook
	pair, _ := generateTLS()

	mux := http.NewServeMux()
	mux.HandleFunc(webhookPathMutate, serveHandler)

	// create or update the mutatingwebhookconfiguration
	err = createOrUpdateMutatingWebhookConfiguration(webhookServiceName, webhookNamespace, k)
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
		},
	}

	// start the HTTP server
	go func() {
		infoLogger.Printf("Starting webhook server on %s", s.Addr)
		if err := s.ListenAndServeTLS("", ""); err != nil {
			log.Fatalf("Failed to start webhook server: %v", err)
		}
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	infoLogger.Printf("Got OS shutdown signal, shutting down webhook server gracefully...")
	s.Shutdown(context.Background()) //nolint:errcheck
}
