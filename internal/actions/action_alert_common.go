package actions

import (
	s "github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	log "github.com/sirupsen/logrus"
)

func newAlertSender(url string) (*router.ServiceRouter, error) {
	sender, err := s.CreateSender(url)
	if err != nil {
		return nil, err
	}

	sender.SetLogger(log.StandardLogger())

	return sender, nil
}
