package actions

import (
	"context"
	"fmt"
	"strings"

	s "github.com/containrrr/shoutrrr"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

var (
	_ ActionInterface                            = &alertDiscord{}
	_ models.AlertInterface[models.AlertDiscord] = &alertDiscord{}
)

type (
	// alertDiscord is an action that sends an alert to a Discord channel.
	alertDiscord struct {
		action
		models.AlertDiscord
	}
)

func init() {
	register(AlertDiscord, &alertDiscord{})
}

// Execute sends the alert message to the Discord channel.
func (a *alertDiscord) Execute(ctx context.Context) error {
	alertConfig, err := a.k.GetValueOrValueFrom(ctx, a.image.Namespace, a.data)
	if err != nil {
		return err
	}

	aC, ok := alertConfig.(v1alpha1.AlertConfig)
	if ok {
		a.AlertDiscord = models.AlertDiscord{
			AlertConfig: aC,
		}
	} else {
		return fmt.Errorf("invalid alert configuration")
	}

	token, webhookID, err := a.parseWebhookURL(ctx)
	if err != nil {
		return err
	}

	sender, err := s.CreateSender(fmt.Sprintf("discord://%s@%s?splitlines=no", token, webhookID))
	if err != nil {
		return fmt.Errorf("error creating Discord sender: %v", err)
	}

	message, err := a.Render()
	if err != nil {
		return fmt.Errorf("error rendering alert message: %v", err)
	}

	var bigErr error
	if errS := sender.Send(message, nil); errS != nil {
		for _, e := range errS {
			if e != nil {
				bigErr = fmt.Errorf("%v: %v", bigErr, e)
			}
		}

		return bigErr
	}

	log.WithField("action", a.GetName()).Info("Alert sent to Discord")

	return nil
}

// ConfigValidation initializes the alertDiscord action.
func (a *alertDiscord) ConfigValidation() error {
	return nil
}

// Render renders the alert message with the provided data.
func (a *alertDiscord) Render() (string, error) {
	aT := alertTemplate[models.AlertDiscord]{
		templateBody:   a.Spec.Discord.TemplateBody,
		tags:           a.tags,
		Image:          *a.action.image,
		AlertInterface: a,
	}
	return aT.Render()
}

// GetSpec returns the alert configuration.
func (a *alertDiscord) GetSpec() models.AlertDiscord {
	return a.AlertDiscord
}

// parsewebhookURL parses the webhook URL from the alert configuration.
// return token and webhookID
func (a *alertDiscord) parseWebhookURL(ctx context.Context) (token, webhookid string, err error) {
	// URL format is https://discord.com/api/webhooks/webhookid/token
	// We need to extract the webhookid and token from the URL

	// Split the URL by /
	// The last element is the token
	// The second to last element is the webhookid

	if a.Spec.Discord == nil {
		return "", "", fmt.Errorf("discord configuration is empty")
	}

	url, err := a.k.GetValueOrValueFrom(ctx, a.Namespace, a.Spec.Discord.WebhookURL)
	if err != nil {
		return "", "", err
	}

	if u, ok := url.(string); ok {
		x := strings.Split(u, "/")

		if len(x) < 2 {
			return "", "", fmt.Errorf("invalid webhook URL")
		}

		return x[len(x)-1], x[len(x)-2], nil
	}

	return "", "", fmt.Errorf("invalid webhook URL")
}

// GetName returns the name of the action.
func (a *alertDiscord) GetName() models.ActionName {
	return AlertDiscord
}
