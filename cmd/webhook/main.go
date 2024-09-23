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
	admissionWebhookAnnotationBase = "kimup.cloudavenue.io"
	// webhookURL 				   = "https://kimup-webhook-mutating.kube-image-updater.svc:8443/mutate" // service
	// pathSpecImage                  = "/spec/containers/image"
)

var (
	webhookURL    string = "https://kimup-webhook-mutating.default.svc:8443/mutate" // outside cluster
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	// 	port                                 int
	// 	envConfigFile                        string
	webhookNamespace   string = "default"
	webhookServiceName string
	webhookConfigName  = "webhookconfig"
	webhookPath        = "/mutate"

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
	infoLogger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// webhook server running namespace
	if os.Getenv("POD_NAMESPACE") != "" {
		webhookNamespace = os.Getenv("POD_NAMESPACE")
	}
}

// Start http server for webhook
func main() {
	var err error
	flag.StringVar(&webhookServiceName, "service-name", "kimup-webhook-mutating", "Kimup Webhook Mutating service name.")

	flag.Parse()

	k, err := client.New("/Users/micheneaudavid/.kube/config")
	if err != nil {
		panic(err)
	}

	// generate cert for webhook
	pair, caPEM := generateTLS()

	mux := http.NewServeMux()
	mux.HandleFunc(webhookPath, serveHandler)

	// create or update the mutatingwebhookconfiguration
	err = createOrUpdateMutatingWebhookConfiguration(caPEM, webhookServiceName, webhookNamespace, k)
	if err != nil {
		errorLogger.Fatalf("Failed to create or update the mutating webhook configuration: %v", err)
	}

	// define http server and server handler
	s := &http.Server{
		Addr:        ":8443",
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{pair},
			MinVersion:   tls.VersionTLS12,
			// InsecureSkipVerify: true, //nolint:gosec
		},
	}

	// start the server
	go func() {
		infoLogger.Printf("Starting webhook server on %s", s.Addr)
		// start TLS server
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
