package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	admissionv1 "k8s.io/api/admission/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"

	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
)

const (
	admissionWebhookAnnotationBase = "kimup.io"
	// pathSpecImage                  = "/spec/containers/image"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	// 	port                                 int
	// 	envConfigFile                        string
	webhookNamespace, webhookServiceName string
	webhookConfigName                    = "webhookconfig"
	webhookPath                          = "/mutate"

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
	webhookNamespace = os.Getenv("POD_NAMESPACE")
}

// Start http server for webhook
func main() {
	var err error
	flag.StringVar(&webhookServiceName, "service-name", "kimup-webhook-mutating", "Kimup Webhook Mutating service name.")

	// generate cert for webhook
	pair, caPEM := generateTLS()

	mux := http.NewServeMux()
	mux.HandleFunc(webhookPath, serveHandler)

	// create or update the mutatingwebhookconfiguration
	err = createOrUpdateMutatingWebhookConfiguration(caPEM, webhookServiceName, webhookNamespace)
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
		},
	}

	// start the server
	go func() {
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

func serveHandler(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		warningLogger.Println("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		warningLogger.Printf("Content-Type=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *admissionv1.AdmissionResponse
	ar := admissionv1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		warningLogger.Printf("Can't decode body: %v", err)
		admissionResponse = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		admissionResponse = mutate(&ar)
	}

	admissionReview := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
	}
	if admissionResponse != nil {
		admissionReview.Response = admissionResponse
		if ar.Request != nil {
			admissionReview.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		warningLogger.Printf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	infoLogger.Printf("Ready to write response ...")
	if _, err := w.Write(resp); err != nil {
		warningLogger.Printf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

// func endpoint /mutate handler
// func mutateHandler(w http.ResponseWriter, r *http.Request) {
// 	// read the body / request
// 	body, err := io.ReadAll(r.Body)
// 	defer r.Body.Close()
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprintf(w, "%s", err)
// 	}

// 	// mutate the request
// 	mutated, err := m.Mutate(body, true)
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprintf(w, "%s", err)
// 	}

// 	// and write it back
// 	w.WriteHeader(http.StatusOK)
// 	_, err = w.Write(mutated)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func generate cert for webhook
func generateTLS() (tls.Certificate, *bytes.Buffer) {
	// generate dns names
	dnsNames := []string{
		webhookServiceName,
		webhookServiceName + "." + webhookNamespace,
		webhookServiceName + "." + webhookNamespace + ".svc",
	}
	commonName := webhookServiceName + "." + webhookNamespace + ".svc"

	caPEM, certPEM, certKeyPEM, err := generateCert([]string{admissionWebhookAnnotationBase}, dnsNames, commonName)
	if err != nil {
		errorLogger.Fatalf("Failed to generate ca and certificate key pair: %v", err)
	}

	pair, err := tls.X509KeyPair(certPEM.Bytes(), certKeyPEM.Bytes())
	if err != nil {
		errorLogger.Fatalf("Failed to load certificate key pair: %v", err)
	}
	return pair, caPEM
}

func generateCert(orgs, dnsNames []string, commonName string) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
	// init CA config
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(2022),
		Subject:               pkix.Name{Organization: orgs},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // expired in 1 year
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// generate private key for CA
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	// create the CA certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// CA certificate with PEM encoded
	caPEM := new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	// new certificate config
	newCert := &x509.Certificate{
		DNSNames:     dnsNames,
		SerialNumber: big.NewInt(1024),
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: orgs,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // expired in 1 year
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// generate new private key
	newPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	// sign the new certificate
	newCertBytes, err := x509.CreateCertificate(rand.Reader, newCert, ca, &newPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// new certificate with PEM encoded
	newCertPEM := new(bytes.Buffer)
	_ = pem.Encode(newCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCertBytes,
	})

	// new private key with PEM encoded
	newPrivateKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(newPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(newPrivateKey),
	})

	return caPEM, newCertPEM, newPrivateKeyPEM, nil
}

// Webhook logic
func createOrUpdateMutatingWebhookConfiguration(caPEM *bytes.Buffer, webhookService, webhookNamespace string) error {
	infoLogger.Println("Initializing the kube client...")

	// create the kube client
	// kubernetes golang libraru provide flag "kubeconfig" to specify the path to the kubeconfig file
	k, err := kubeclient.New(flag.Lookup("kubeconfig").Value.String())
	if err != nil {
		panic(err)
	}

	// kubeconfig := os.Getenv("KUBECONFIG")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	return err
	// }
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	mutatingWebhookConfigV1Client := clientset.AdmissionregistrationV1()

	infoLogger.Printf("Creating or updating the mutatingwebhookconfiguration: %s", webhookConfigName)
	fail := admissionregistrationv1.Fail
	sideEffect := admissionregistrationv1.SideEffectClassNone
	mutatingWebhookConfig := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: webhookConfigName,
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{{
			Name:                    admissionWebhookAnnotationBase,
			AdmissionReviewVersions: []string{"v1", "v1beta1"},
			SideEffects:             &sideEffect,
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				CABundle: caPEM.Bytes(), // self-generated CA for the webhook
				Service: &admissionregistrationv1.ServiceReference{
					Name:      webhookService,
					Namespace: webhookNamespace,
					Path:      &webhookPath,
				},
			},
			Rules: []admissionregistrationv1.RuleWithOperations{
				{
					Operations: []admissionregistrationv1.OperationType{
						admissionregistrationv1.Update,
					},
					Rule: admissionregistrationv1.Rule{
						APIGroups:   []string{""},
						APIVersions: []string{"v1"},
						Resources:   []string{"pods"},
					},
				},
			},
			// NamespaceSelector: &metav1.LabelSelector{
			// 	MatchLabels: map[string]string{
			// 		"env-injection": "enabled",
			// 	},
			// },
			FailurePolicy: &fail,
		}},
	}

	foundWebhookConfig, err := mutatingWebhookConfigV1Client.MutatingWebhookConfigurations().Get(context.TODO(), webhookConfigName, metav1.GetOptions{})
	switch {
	case err != nil && apierrors.IsNotFound(err):
		if _, err := mutatingWebhookConfigV1Client.MutatingWebhookConfigurations().Create(context.TODO(), mutatingWebhookConfig, metav1.CreateOptions{}); err != nil {
			warningLogger.Printf("Failed to create the mutatingwebhookconfiguration: %s", webhookConfigName)
			return err
		}
		infoLogger.Printf("Created mutatingwebhookconfiguration: %s", webhookConfigName)
	case err != nil:
		warningLogger.Printf("Failed to check the mutatingwebhookconfiguration: %s", webhookConfigName)
		return err
	default:
		// there is an existing mutatingWebhookConfiguration
		if len(foundWebhookConfig.Webhooks) != len(mutatingWebhookConfig.Webhooks) ||
			!(foundWebhookConfig.Webhooks[0].Name == mutatingWebhookConfig.Webhooks[0].Name &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].AdmissionReviewVersions, mutatingWebhookConfig.Webhooks[0].AdmissionReviewVersions) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].SideEffects, mutatingWebhookConfig.Webhooks[0].SideEffects) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].FailurePolicy, mutatingWebhookConfig.Webhooks[0].FailurePolicy) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].Rules, mutatingWebhookConfig.Webhooks[0].Rules) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].NamespaceSelector, mutatingWebhookConfig.Webhooks[0].NamespaceSelector) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.CABundle, mutatingWebhookConfig.Webhooks[0].ClientConfig.CABundle) &&
				reflect.DeepEqual(foundWebhookConfig.Webhooks[0].ClientConfig.Service, mutatingWebhookConfig.Webhooks[0].ClientConfig.Service)) {
			mutatingWebhookConfig.ObjectMeta.ResourceVersion = foundWebhookConfig.ObjectMeta.ResourceVersion
			if _, err := mutatingWebhookConfigV1Client.MutatingWebhookConfigurations().Update(context.TODO(), mutatingWebhookConfig, metav1.UpdateOptions{}); err != nil {
				warningLogger.Printf("Failed to update the mutatingwebhookconfiguration: %s", webhookConfigName)
				return err
			}
			infoLogger.Printf("Updated the mutatingwebhookconfiguration: %s", webhookConfigName)
		}
		infoLogger.Printf("The mutatingwebhookconfiguration: %s already exists and has no change", webhookConfigName)
	}

	return nil
}

// func mutate the request
func mutate(ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		warningLogger.Printf("Could not unmarshal raw object: %v", err)
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	infoLogger.Printf("AdmissionReview for Kind=%v, Namespace=%v Name=%v (%v) UID=%v patchOperation=%v UserInfo=%v",
		req.Kind, req.Namespace, req.Name, pod.Name, req.UID, req.Operation, req.UserInfo)

	// annotations := map[string]string{admissionWebhookAnnotationStatusKey: "updated"}
	patchBytes, err := createPatch(&pod)
	if err != nil {
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}
	infoLogger.Printf("AdmissionResponse: patch=%v\n", string(patchBytes))
	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// create mutation patch for resoures.
func createPatch(pod *corev1.Pod) ([]byte, error) {
	var patch []patchOperation

	// patch = append(patch, addEnvironnement(pod.Spec.Containers, config.Env, "/spec/containers")...)
	patch = append(patch, updateImage(pod)...)

	return json.Marshal(patch)
}

// updateImage updates the image of the container
func updateImage(pod *corev1.Pod) (patch []patchOperation) {
	// Handle Containers
	for i := range pod.Spec.Containers {
		container := &pod.Spec.Containers[i]
		patch = append(patch, patchOperation{
			Op:    "replace",
			Path:  container.Image,
			Value: "debian:1.2.3",
		})
	}
	return patch
}
