package triggers

import (
	"github.com/gookit/event"
	"github.com/sirupsen/logrus"

	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

type (
	Name      string
	EventName string
)

const (
	RefreshImage  EventName = "refresh.image"
	RefreshStatus EventName = "refresh.status"

	Crontab Name = "crontab"
	Webhook Name = "webhook"
)

func (e EventName) String() string {
	return string(e)
}

func Trigger(e EventName, namespace, imageName string) (event.Event, error) {
	log.
		WithFields(logrus.Fields{
			"namespace": namespace,
			"image":     imageName,
		}).Infof("Triggering event %s", e.String())

	event.Async(e.String(), event.M{"namespace": namespace, "image": imageName})
	return nil, nil
}
