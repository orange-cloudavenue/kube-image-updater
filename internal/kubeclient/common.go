package kubeclient

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type UnstructuredFunc interface {
	UnstructuredContent() map[string]interface{}
}

func decodeUnstructured[T any](v UnstructuredFunc) (t T, err error) {
	if err := runtime.DefaultUnstructuredConverter.
		FromUnstructured(v.UnstructuredContent(), &t); err != nil {
		return t, fmt.Errorf("failed to convert resource: %w", err)
	}

	return
}

func DecodeUnstructured[T any](v *unstructured.Unstructured) (t T, err error) {
	return decodeUnstructured[T](v)
}

func encodeUnstructured[T any](t T) (*unstructured.Unstructured, error) {
	x, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&t)
	if err != nil {
		return nil, fmt.Errorf("failed to convert resource: %w", err)
	}

	return &unstructured.Unstructured{Object: x}, nil
}

func EncodeUnstructured[T any](t T) (*unstructured.Unstructured, error) {
	return encodeUnstructured[T](t)
}
