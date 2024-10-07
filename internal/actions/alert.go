package actions

import (
	"bytes"
	"html/template"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

const (
	defaultAlertTemplate = `
	Kimup alert for image update:
	{{ .Namespace }}/{{ .Name }}

	Image **{{ .ImageName }}:{{ .ActualTag }}** has a new tag available: **{{ .NewTag }}**

	Available tags:
{{ range .AvailableTags -}}
	- {{ . }}
{{ end }}
	`
)

type (
	alertTemplate[T any] struct {
		templateBody string
		tags         models.Tags
		v1alpha1.Image
		models.AlertInterface[T]
	}

	alertTemplateData struct {
		// * Common
		Namespace string
		Name      string
		ImageName string
		BaseTag   string

		// * Tags
		NewTag        string
		ActualTag     string
		AvailableTags []string
	}
)

func (a *alertTemplate[T]) Render() (string, error) {
	if a.templateBody == "" {
		a.templateBody = defaultAlertTemplate
	}

	t, err := template.New("alert").Parse(a.templateBody)
	if err != nil {
		return "", err
	}

	data := alertTemplateData{
		Namespace: a.Namespace,
		Name:      a.Name,
		ImageName: a.Image.Spec.Image,
		BaseTag:   a.Image.Spec.BaseTag,

		NewTag:        a.tags.New,
		ActualTag:     a.tags.Actual,
		AvailableTags: a.tags.AvailableTags,
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
