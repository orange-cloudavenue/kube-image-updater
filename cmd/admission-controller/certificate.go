package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	v1 "k8s.io/client-go/applyconfigurations/admissionregistration/v1"

	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
)

// generateTLS generates a self-signed certificate for the webhook server
// and returns the certificate and the CA certificate
// The certificate is generated with the following DNS names:
// - webhookServiceName
// - webhookServiceName.webhookNamespace
// - webhookServiceName.webhookNamespace.svc
func generateTLS() (keyPair tls.Certificate, caPEM *bytes.Buffer, err error) {
	// generate dns names
	dnsNames := []string{
		webhookServiceName,
		webhookServiceName + "." + webhookNamespace,
		webhookServiceName + "." + webhookNamespace + ".svc",
		// webhookServiceName + "." + webhookNamespace + ".svc" + ".cluster.local",
	}
	commonName := webhookServiceName + "." + webhookNamespace + ".svc"

	caPEM, certPEM, certKeyPEM, err := generateCert([]string{webhookBase}, dnsNames, commonName)
	if err != nil {
		return
	}

	keyPair, err = tls.X509KeyPair(certPEM.Bytes(), certKeyPEM.Bytes())
	if err != nil {
		return
	}
	return
}

// generateCert generates a self-signed certificate with the given organizations, DNS names, and common name
// The certificate is valid for 1 year
// The certificate is signed by the CA certificate
// The CA certificate is generated with the given organizations
// it resurns the CA, certificate and private key in PEM format.
func generateCert(orgs, dnsNames []string, commonName string) (caPEM, newCertPEM, newPrivateKeyPEM *bytes.Buffer, err error) {
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
	caPEM = new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// print CA certificate if insideCluster is false
	if !insideCluster {
		writeNewCA(caPEM, manifestWebhookPath)
		time.Sleep(2 * time.Second)
		applyManifest(manifestWebhookPath)
	}

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
	newCertPEM = new(bytes.Buffer)
	err = pem.Encode(newCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: newCertBytes,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// new private key with PEM encoded
	newPrivateKeyPEM = new(bytes.Buffer)
	err = pem.Encode(newPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(newPrivateKey),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return caPEM, newCertPEM, newPrivateKeyPEM, nil
}

func writeNewCA(caPEM *bytes.Buffer, filePath string) {
	newCABundle := base64.StdEncoding.EncodeToString(caPEM.Bytes())

	// Lire le fichier
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "caBundle:") {
			line = "      caBundle: " + "\"" + newCABundle + "\""
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return
	}

	// Écrire les modifications dans le fichier
	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return
		}
	}
	writer.Flush()
}

func applyManifest(file string) {
	// read the manifest file
	manifestBytes, err := os.ReadFile(file)
	if err != nil {
		return
	}

	// decode the manifest to unstructured object
	decoder := serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	obj := &unstructured.Unstructured{}
	_, _, err = decoder.Decode(manifestBytes, nil, obj)
	if err != nil {
		return
	}

	// convert the unstructured object to typed object
	mutatingWebhookConfiguration, err := kubeclient.DecodeUnstructured[v1.MutatingWebhookConfigurationApplyConfiguration](obj)
	if err != nil {
		return
	}

	// apply the manifest
	if _, err := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Apply(
		context.TODO(),
		&mutatingWebhookConfiguration,
		metav1.ApplyOptions{Force: true, FieldManager: "kumi-webhook"},
	); err != nil {
		return
	}
}
