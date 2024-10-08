---
hide:
  - toc
---

# Discord Alert

The `discord` alert allows you to send a message to a Discord channel when an image is updated.

## Who to use

!!! warning "Require AlertConfig"
    The `discord` alert requires an `AlertConfig` resource to be created.

Create an `AlertConfig` resource with the `discord` alert.

### Setting

The CRD schema for the `AlertConfig` resource is available on [doc.crds.dev](https://doc.crds.dev/github.com/orange-cloudavenue/kube-image-updater@latest)

#### Examples

**1 - Create kubernetes secret**

```bash
kubectl create secret generic discord-secret --from-literal=webhookURL=https://discord.com/api/webhooks/1234567890/ABCDEFGHIJKLMN --dry-run=client -o yaml > discord-secret.yaml

kubectl apply -f discord-secret.yaml
```

**2 - Create AlertConfig**

```yaml hl_lines="6-13"
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: AlertConfig
metadata:
  name: demo
spec:
  discord:
    webhookURL: # (1)
      valueFrom: # (2)
        secretKeyRef:
          name: discord-secret
          key: webhookURL
    templateBody: | # (3)
      New dev version {{ .NewTag }} is available for {{ .ImageName }}.
```

1. The `webhookURL` is the URL of the Discord webhook. [How to create a Discord webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
2. The `valueFrom` field allows you to reference a secret key. The `discord-secret` secret must be created with a `webhookURL` key.
3. The `templateBody` field is the custom message. For more information about the template, you can check the [template documentation](getting-start.md#template-body-alert-message).

**3 - Create Image**

In this example, if a dev version is detected, an alert will be sent to the Discord channel and the new image tag will be applied.

```yaml hl_lines="18-22"
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Image
metadata:
  labels:
    app.kubernetes.io/managed-by: kustomize
  name: demo
spec:
  image: registry.127.0.0.1.nip.io/demo
  baseTag: v0.0.4
  triggers:
    - [...]
  rules:
    - name: Automatic apply on dev version
      type: regex
      # Match v1.2.3-dev1 version
      value: "^v?[0-9].[0-9].[0-9]-dev[0-9]$"
      actions:
        - type: alert-discord
          data:
            valueFrom:
              alertConfigRef: # (1)
                name: demo
        - type: apply
```

1. The `alertConfigRef` field allows you to reference the `AlertConfig` resource in the same namespace.

## Fields

| Field | Description | Mandatory | Default |
|-------|-------------|:-----------:|---------|
| `webhookURL` | The URL of the Discord webhook. | :white_check_mark: | |
| `templateBody` | The custom message. | | [Default message](getting-start.md#template-body-alert-message) |
| `templateFile` | The path to the custom message file. | | |
