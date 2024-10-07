---
hide:
  - toc
---

# Getting started Alert

Kimup provide multiple actions to alert you when an image is updated.

!!! info "new actions are coming soon"
    We are working on new actions to alert you when an image is updated. Stay tuned!
    [Github issue](https://github.com/orange-cloudavenue/kube-image-updater/issues?q=sort:updated-desc+is:issue+is:open+label:action-alert)

List of available alerts :

* [Discord](discord.md)
* Email (coming soon)
* Slack (coming soon)
* Webhook (coming soon)

## Advanced usage

### Template body alert message

Alert have a custom body template to customize the message sent.

The template is a [Go template](https://pkg.go.dev/text/template) with the following variables:

| Variable | Description | Type | Example |
| :--- | :--- | :---: | :--- |
| `{{ .Namespace }}` | The namespace of the resource | string |`default` |
| `{{ .Name }}` | The name of the resource | string | `demo` |
| `{{ .ImageName }}` | The image name | string | `ghcr.io/orange-cloudavenue/kube-image-updater` |
| `{{ .BaseTag }}` | The base tag | string | `v0.0.19` |
| `{{ .NewTag }}` | The new tag | string | `v0.0.22` |
| `{{ .ActualTag }}` | The actual tag | string | `v0.0.21` |
| `{{ .AvailableTags }}` | The available tags | slice | `v0.0.19, v0.0.20, v0.0.21, v0.0.22` |

**Default template body alert message**

```
	Kimup alert for image update:
	{{ .Namespace }}/{{ .Name }}

	Image **{{ .ImageName }}:{{ .ActualTag }}** has a new tag available: **{{ .NewTag }}**

	Available tags:
{{ range .AvailableTags -}}
	- {{ . }}
{{ end }}
```
