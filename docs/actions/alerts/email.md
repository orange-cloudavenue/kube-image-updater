---
hide:
  - toc
---

# Email Alert

The `email` alert allows you to send an email when an image is updated.

## Who to use

!!! warning "Require AlertConfig"
    The `email` alert requires an `AlertConfig` resource to be created.

Create an `AlertConfig` resource with the `email` alert.

### Setting

The CRD schema for the `AlertConfig` resource is available on [doc.crds.dev](https://doc.crds.dev/github.com/orange-cloudavenue/kube-image-updater@latest)

**1 - Create kubernetes secret**

```bash
kubectl create secret generic email-secret \
    --from-literal=smtpHost=smtp.example.com\
    --from-literal=smtpPort=587 \
    --from-literal=smtpUsername=foo \
    --from-literal=smtpPassword=bar \
    --dry-run=client -o yaml > email-secret.yaml

kubectl apply -f email-secret.yaml
```

**2 - Create AlertConfig**
<!-- Loaded from file because the vars in template are rendered by mkdocs-macros plugins and generate a error -->
```yaml hl_lines="6-34"
--8<-- "docs/actions/alerts/email-alert-config.yaml"
```

1. The `host` is the SMTP server host.
2. The `port` is the SMTP server port. Default value is `25`.
3. The `username` is the SMTP server username.
4. The `password` is the SMTP server password.
5. The `toAddress` is the list of email addresses to send the alert.

**3 - Create Image**

In this example, if a dev version is detected, an alert will be sent to the email address and the new image tag will be applied.

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
        - type: alert-email
          data:
            valueFrom:
              alertConfigRef: # (1)
                name: demo
        - type: apply
```

1. The `alertConfigRef` field allows you to reference the `AlertConfig` resource in the same namespace.

## Fields

See the list of fields available for the `email` alert in the [doc.crds.dev](https://doc.crds.dev/github.com/orange-cloudavenue/kube-image-updater/kimup.cloudavenue.io/AlertConfig/v1alpha1@{{git.short_tag}}#spec-email)
