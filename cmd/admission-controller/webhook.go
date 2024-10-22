package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
	"github.com/orange-cloudavenue/kube-image-updater/internal/patch"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

// func serveHandler
func ServeHandler(w http.ResponseWriter, r *http.Request) {
	// Prometheus metrics
	metrics.AdmissionController().RequestTotal.Inc()
	timeAC := metrics.AdmissionController().RequestDuration.NewTimer()
	defer timeAC.ObserveDuration()

	var body []byte
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		// increment the total number of errors
		metrics.AdmissionController().RequestErrorTotal.Inc()

		log.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		// increment the total number of errors
		metrics.AdmissionController().RequestErrorTotal.Inc()

		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *admissionv1.AdmissionResponse
	ar := admissionv1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		// increment the total number of errors
		metrics.AdmissionController().RequestErrorTotal.Inc()

		log.WithError(err).Warn("Can't decode body")
		admissionResponse = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		admissionResponse = mutate(r.Context(), &ar)
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
		// increment the total number of errors
		metrics.AdmissionController().RequestErrorTotal.Inc()

		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		// increment the total number of errors
		metrics.AdmissionController().RequestErrorTotal.Inc()

		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

// func mutate the request
func mutate(ctx context.Context, ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	req := ar.Request
	var pod corev1.Pod
	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	log.WithFields(logrus.Fields{
		"Kind":      req.Kind,
		"Namespace": req.Namespace,
		"Name":      req.Name,
		"UID":       req.UID,
		"Operation": req.Operation,
		"UserInfo":  req.UserInfo,
	}).Info("AdmissionReview")

	// create patch
	patchBytes, err := createPatch(ctx, &pod)
	if err != nil {
		return &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}
	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// create mutation patch for pod.
func createPatch(ctx context.Context, pod *corev1.Pod) ([]byte, error) {
	// Metrics - increment the total number of patch
	metrics.AdmissionController().PatchTotal.Inc()
	timePatch := metrics.AdmissionController().PatchDuration.NewTimer()
	defer timePatch.ObserveDuration()

	var err error
	// find annotation enabled
	an := annotations.New(ctx, pod)
	if !an.Enabled().Get() {
		// increment the total number of errors
		metrics.AdmissionController().PatchErrorTotal.Inc()

		return nil, fmt.Errorf("annotation not enabled")
	}

	// var patch []patchOperation
	p := patch.NewBuilder()

	log.
		WithFields(logrus.Fields{
			"Namespace": pod.Namespace,
			"Name":      pod.Name,
		}).Info("Generate Patch")

	for i, container := range pod.Spec.Containers {
		imageP := utils.ImageParser(container.Image)

		// TODO Why is this not used? Annotation is never set.
		crdName, _ := an.Images().Get(imageP.GetImageWithoutTag())

		// If crdName is empty, it means that we need to find it
		var image v1alpha1.Image
		if crdName == "" {
			// find the image associated with the pod
			image, err = kubeClient.Image().Find(ctx, pod.Namespace, imageP.GetImageWithoutTag())
			if err != nil {
				// increment the total number of errors
				metrics.AdmissionController().PatchErrorTotal.Inc()

				log.
					WithFields(logrus.Fields{
						"Namespace":      pod.Namespace,
						"Name":           pod.Name,
						"Container":      container.Name,
						"ContainerImage": imageP.GetImageWithoutTag(),
					}).
					WithError(err).Error("Failed to find kind Image")
				continue
			}
		} else {
			image, err = kubeClient.Image().Get(ctx, pod.Namespace, crdName)
			if err != nil {
				// increment the total number of errors
				metrics.AdmissionController().PatchErrorTotal.Inc()

				log.
					WithFields(logrus.Fields{
						"Namespace":      pod.Namespace,
						"Name":           pod.Name,
						"Container":      container.Name,
						"ContainerImage": crdName,
					}).WithError(err).Error("Failed to get kind Image")
				continue
			}
		}

		// Set the image to the pod
		if image.ImageIsEqual(container.Image) {
			p.AddPatch(patch.OpReplace, fmt.Sprintf("/spec/containers/%d/image", i), image.GetImageWithTag())
		}

		// Annotations
		an.Containers().Set(container.Name, image.Name)
	}

	// update the annotation
	p.AddRawPatches(an.Containers().BuildPatches())

	return p.Generate()
}
