---
hide:
  - toc
---

# Custom Resource Definition `AlertConfig`

This is a custom resource definition for an alert configuration. It is used to setting up alerts for the image update.
`AlertConfig` is a namespaced resource.

## Basic example

```yaml
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: AlertConfig
metadata:
  name: demo
spec:
  discord:
    webhookURL:
      valueFrom:
        secretKeyRef:
          name: discord-secret
          key: webhookURL
```

Each alert type has its own configuration. The `AlertConfig` resource defines the configuration for the alerts.
