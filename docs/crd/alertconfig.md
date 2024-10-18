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

Each alert type has its own configuration.
The `AlertConfig` resource defines the configuration for the alerts.

## Configuration

Kimup Operator uses a dedicated kimup CRD to create and manage AlertConfig resources. The CRD allows various configurations to define the behaviour of the image. See [docs.crds.dev](https://doc.crds.dev/github.com/orange-cloudavenue/kube-image-updater/kimup.cloudavenue.io/AlertConfig/v1alpha1) for more information about the AlertConfig CRD.
