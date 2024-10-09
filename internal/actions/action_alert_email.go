package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/containrrr/shoutrrr/pkg/types"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

var (
	_ models.ActionInterface                   = &alertEmail{}
	_ models.AlertInterface[models.AlertEmail] = &alertEmail{}
)

type (
	// alertEmail is an action that sends an alert via email.
	alertEmail struct {
		action
		models.AlertEmail
	}
)

func init() {
	register(AlertEmail, &alertEmail{})
}

// Execute sends the alert message via email.
func (a *alertEmail) Execute(ctx context.Context) error {
	alertConfig, err := a.k.GetValueOrValueFrom(ctx, a.image.Namespace, a.data)
	if err != nil {
		return err
	}

	aC, ok := alertConfig.(v1alpha1.AlertConfig)
	if !ok {
		return fmt.Errorf("invalid alert configuration")
	}

	a.AlertEmail = models.AlertEmail{
		AlertConfig: aC,
	}

	// Construct URL
	url, err := a.constructURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to construct URL: %w", err)
	}

	// Create sender
	sender, err := newAlertSender(url)
	if err != nil {
		return fmt.Errorf("failed to create sender: %w", err)
	}

	// Render message
	message, err := a.Render()
	if err != nil {
		return fmt.Errorf("failed to render message: %w", err)
	}

	var bigErr error
	if errS := sender.Send(message, &types.Params{}); errS != nil {
		for _, e := range errS {
			if e != nil {
				bigErr = fmt.Errorf("failed to send email: %w", e)
			}
		}

		return bigErr
	}

	return nil
}

// ConfigValidation validates the alert email configuration.
func (a *alertEmail) ConfigValidation() error {
	return nil
}

// Render renders the alert email message.
func (a *alertEmail) Render() (string, error) {
	aT := alertTemplate[models.AlertEmail]{
		templateBody:   a.Spec.Email.TemplateBody,
		tags:           a.tags,
		Image:          *a.action.image,
		AlertInterface: a,
	}
	return aT.Render()
}

// GetSpec returns the alert email spec.
func (a *alertEmail) GetSpec() models.AlertEmail {
	return a.AlertEmail
}

// GetName returns the name of the action.
func (a *alertEmail) GetName() models.ActionName {
	return AlertEmail
}

// construct url
func (a *alertEmail) constructURL(ctx context.Context) (string, error) {
	url := "smtp://"

	// * Username
	if username, err := a.k.GetValueOrValueFrom(ctx, a.GetNamespace(), a.Spec.Email.Username); err == nil {
		url += username.(string)
		// * Password
		if password, err := a.k.GetValueOrValueFrom(ctx, a.GetNamespace(), a.Spec.Email.Password); err == nil {
			url += ":" + password.(string) + "@"
		}
	}

	// * Host
	if host, err := a.k.GetValueOrValueFrom(ctx, a.GetNamespace(), a.Spec.Email.Host); err == nil && host.(string) != "" {
		url += host.(string)
		// * Port
		if port, err := a.k.GetValueOrValueFrom(ctx, a.GetNamespace(), a.Spec.Email.Port); err == nil {
			url += ":" + port.(string)
		}
	} else {
		return "", fmt.Errorf("failed to get host: %w", err)
	}

	// Query/Param Props
	params := map[string]string{
		"from":       a.Spec.Email.FromAddress,
		"to":         strings.Join(a.Spec.Email.ToAddress, ","),
		"clientHost": a.Spec.Email.ClientHost,
		"encryption": a.Spec.Email.Encryption,
		"fromName":   a.Spec.Email.FromName,
		"auth":       a.Spec.Email.Auth,
		"useHTML":    "no",
	}

	// UseHTML / UseStartTLS
	if a.Spec.Email.UseHTML {
		params["useHTML"] = "yes"
	}

	if a.Spec.Email.UseStartTLS {
		params["starttls"] = "yes"
	}

	// Subject
	if a.Spec.Email.TemplateSubject == "" {
		a.Spec.Email.TemplateSubject = "Kimup - New tag is available for {{ .ImageName }}"
	}

	aT := alertTemplate[models.AlertEmail]{
		templateBody:   a.Spec.Email.TemplateSubject,
		tags:           a.tags,
		Image:          *a.action.image,
		AlertInterface: a,
	}

	subject, err := aT.Render()
	if err != nil {
		return "", fmt.Errorf("failed to render subject: %w", err)
	}

	params["subject"] = subject

	// Construct URL
	url += "?"
	for k, v := range params {
		url += k + "=" + v + "&"
	}

	return url, nil
}
