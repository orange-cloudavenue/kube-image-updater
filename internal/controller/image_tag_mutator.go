package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

func (i *ImageTagMutator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mgr.GetWebhookServer().Register("/mutate/image-tag", &webhook.Admission{Handler: i.SetupHandler()})
	return nil
}

func (i *ImageTagMutator) SetupHandler() admission.Handler {
	i.decoder = admission.NewDecoder(i.Scheme)
	return i
}

// +kubebuilder:webhook:path=/mutate/image-tag,mutating=true,failurePolicy=fail,groups="",resources=pods,sideEffects=None,verbs=create;update,versions=v1,name=mutator.kimup.cloudavenue.io,admissionReviewVersions=v1

var _ admission.Handler = &ImageTagMutator{}

// podAnnotator annotates Pods
type ImageTagMutator struct {
	client.Client
	KubeAPIClient *kubeclient.Client
	Scheme        *runtime.Scheme
	decoder       admission.Decoder
}

func (i *ImageTagMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := logf.FromContext(ctx)

	pod := &corev1.Pod{}

	err := i.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	for c, container := range pod.Spec.Containers {
		imageP := utils.ImageParser(container.Image)

		// find the image associated with the pod
		image, err := i.KubeAPIClient.Image().Find(ctx, pod.Namespace, imageP.GetImageWithoutTag())
		if err != nil {
			// increment the total number of errors
			metrics.Mutator().PatchErrorTotal.Inc()

			log.Error(err, "Failed to find kind Image")
			continue
		}

		log.Info(fmt.Sprintf("Mutating container %s with image %s to %s", container.Name, container.Image, image.GetImageWithTag()))

		// Set the image to the pod
		if image.ImageIsEqual(container.Image) {
			pod.Spec.Containers[c].Image = image.GetImageWithTag()
		}
	}

	// Marshal the pod and return a patch response
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		log.Error(err, "Failed to mutate the pod")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}
