package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/metrics"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
)

func (i *ImageTagMutator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	wh := admission.WithCustomDefaulter(mgr.GetScheme(), &corev1.Pod{}, i)
	wh.WithRecoverPanic(true)
	mgr.GetWebhookServer().Register("/mutate/image-tag", &webhook.Admission{Handler: wh})
	return nil
}

// +kubebuilder:webhook:path=/mutate/image-tag,mutating=true,failurePolicy=fail,groups="",resources=pods,sideEffects=None,verbs=create;update,versions=v1,name=mutator.kimup.cloudavenue.io,admissionReviewVersions=v1

var _ admission.CustomDefaulter = &ImageTagMutator{}

// podAnnotator annotates Pods
type ImageTagMutator struct {
	client.Client
	KubeAPIClient *kubeclient.Client
}

func (i *ImageTagMutator) Default(ctx context.Context, obj runtime.Object) error {
	log := logf.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	an := annotations.New(ctx, pod)
	if !an.Enabled().Get() {
		// increment the total number of errors
		metrics.AdmissionController().PatchErrorTotal.Inc()
		log.Info(fmt.Sprintf("annotation not enabled for pod %s/%s. Ignore it", pod.Namespace, pod.Name))
		// Return nil because we don't want to mutate the pod
		return nil
	}

	for _, container := range pod.Spec.Containers {
		imageP := utils.ImageParser(container.Image)

		// TODO Why is this not used? Annotation is never set.
		crdName, _ := an.Images().Get(imageP.GetImageWithoutTag())

		// If crdName is empty, it means that we need to find it
		var (
			image v1alpha1.Image
			err   error
		)

		if crdName == "" {
			// find the image associated with the pod
			image, err = i.KubeAPIClient.Image().Find(ctx, pod.Namespace, imageP.GetImageWithoutTag())
			if err != nil {
				// increment the total number of errors
				metrics.AdmissionController().PatchErrorTotal.Inc()

				log.Error(err, "Failed to find kind Image")
				continue
			}
		} else {
			image, err = i.KubeAPIClient.Image().Get(ctx, pod.Namespace, crdName)
			if err != nil {
				// increment the total number of errors
				metrics.AdmissionController().PatchErrorTotal.Inc()

				log.Error(err, "Failed to get kind Image")
				continue
			}
		}

		container.Image = image.GetImageWithTag()
		// // Set the image to the pod
		// if image.ImageIsEqual(container.Image) {
		// }

		// Annotations
		// an.Containers().Set(container.Name, image.Name)
	}

	return nil
}
