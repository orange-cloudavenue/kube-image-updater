package annotations

import (
	"context"
)

type (
	annotations map[string]string

	aChan chan struct {
		key   AnnotationKey
		value string
	}

	Annotation struct {
		ctx context.Context
		annotations
	}

	AnnotationKey string
)

// AnnotationKey is the key used to store the image in the annotation
var (
	KeyAction        AnnotationKey = "kimup.cloudavenue.io" + "/action"
	KeyTag           AnnotationKey = "kimup.cloudavenue.io" + "/tag"
	KeyCheckSum      AnnotationKey = "kimup.cloudavenue.io" + "/checksum"
	KeyEnabled       AnnotationKey = "kimup.cloudavenue.io" + "/enabled"
	KeyFailurePolicy AnnotationKey = "kimup.cloudavenue.io" + "/failure-policy"
)

type (
	KubeAnnotationInterface interface {
		GetAnnotations() map[string]string
	}
)

func New(ctx context.Context, object KubeAnnotationInterface) Annotation {
	return Annotation{
		ctx:         ctx,
		annotations: object.GetAnnotations(),
	}
}

// * Global

func (a *Annotation) Remove(key AnnotationKey) {
	delete(a.annotations, string(key))
}

// * Generic funcs

func (aC aChan) Send(key AnnotationKey, value string) {
	aC <- struct {
		key   AnnotationKey
		value string
	}{key, value}
}
