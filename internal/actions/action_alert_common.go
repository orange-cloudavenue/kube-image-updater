package actions

import (
	s "github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"

	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

func newAlertSender(url string) (*router.ServiceRouter, error) {
	sender, err := s.CreateSender(url)
	if err != nil {
		return nil, err
	}

	sender.SetLogger(log.GetLogger())

	return sender, nil
}
