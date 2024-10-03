package kubeclient

import (
	"k8s.io/apimachinery/pkg/watch"
)

type (
	WatchInterface[T any] struct {
		Type  watch.EventType
		Value T
	}
)
