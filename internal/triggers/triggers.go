package triggers

import (
	"github.com/gookit/event"
	log "github.com/sirupsen/logrus"
)

type EventName string

const (
	RefreshImage EventName = "refresh.image"
)

func (e EventName) String() string {
	return string(e)
}

func Trigger(e EventName, namespace, imageName string) (event.Event, error) {
	log.Infof("Triggering event %s for image %s in namespace %s", e.String(), imageName, namespace)
	err, x := event.Fire(e.String(), event.M{"namespace": namespace, "image": imageName})
	return x, err
}
