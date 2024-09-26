package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
)

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
		warningLogger.Printf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	infoLogger.Printf("Ready to write response ...")
	if _, err := w.Write(resp); err != nil {
		warningLogger.Printf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

// func mutate the request
func mutate(ctx context.Context, ar *admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
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

	// create patch
	patchBytes, err := createPatch(ctx, &pod)
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

// create mutation patch for pod.
func createPatch(ctx context.Context, pod *corev1.Pod) ([]byte, error) {
	var err error
	// find annotation enabled
	an := annotations.New(ctx, pod)
	if !an.Enabled().Get() {
		return nil, fmt.Errorf("annotation not enabled")
	}

	var patch []patchOperation
	infoLogger.Printf("Generate Patch for: %v\n", pod.Name)

	for i, container := range pod.Spec.Containers {
		// check if an annotation exist
		crdName, _ := an.Images().Get(container.Image)

		// If crdName is empty, it means that we need to find it
		var image v1alpha1.Image
		if crdName == "" {
			// find the image associated with the pod
			image, err = kubeClient.FindImage(ctx, pod.Namespace, container.Image)
			if err != nil {
				warningLogger.Printf("failed to find image: %v", err)
				continue
			}
		} else {
			image, err = kubeClient.GetImage(ctx, pod.Namespace, crdName)
			if err != nil {
				warningLogger.Printf("failed to get image: %v", err)
				continue
			}
		}

		// Set the image to the pod
		if image.ImageIsEqual(container.Image) {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  fmt.Sprintf("/spec/containers/%d/image", i),
				Value: image.GetImageWithTag(),
			})
		}

		an.Containers().Set(container.Name, image.GetImageWithTag())
	}

	// update the annotation
	patch = append(patch, patchOperation{
		Op:    "replace",
		Path:  "/metadata/annotations",
		Value: pod.GetAnnotations(),
	})

	// _, err = kubeClient.GetKubeClient().CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
	// if err != nil {
	// 	warningLogger.Printf("failed to update pod annotation: %v", err)
	// }

	// patch = append(patch, updateImage(pod.Spec.Containers)...)
	debugLogger.Printf("Patch created: %v\n", patch)

	return json.Marshal(patch)
}
