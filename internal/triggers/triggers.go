package triggers

import (
	"github.com/gookit/event"
)

type EventName string

const (
	RefreshImage EventName = "refresh.image"
)

func (e EventName) String() string {
	return string(e)
}

func Trigger(e EventName, namespace, imageName string) {
	event.Async(e.String(), event.M{"namespace": namespace, "image": imageName})
}
