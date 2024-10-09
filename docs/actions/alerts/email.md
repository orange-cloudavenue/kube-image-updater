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

```yaml hl_lines="6-34"
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: AlertConfig
metadata:
  name: demo
spec:
    email:
        host: # Required (1)
          valueFrom:
            secretKeyRef:
            name: email-secret
            key: smtpHost
        port: # Optionnal (2)
          valueFrom:
            secretKeyRef:
            name: email-secret
            key: smtpPort
        username: # Optionnal (3)
          valueFrom:
            secretKeyRef:
            name: email-secret
            key: smtpUsername
        password: # Optionnal (4)
          valueFrom:
            secretKeyRef:
            name: email-secret
            key: smtpPassword
        fromAddress: noreply@bar.com
        toAddress: # Required (5)
          - foo@bar.com
          - bar@foo.com
        templateBody: |
          New dev version {{ .NewTag }} is available for {{ .ImageName }}.
        templateSubject: |
          New version available for {{ .ImageName }}
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

| Field            | Description                                                                                       | Mandatory | Default   |
|------------------|---------------------------------------------------------------------------------------------------|:---------:|-----------|
| `host`           | Host specifies the SMTP server to connect to.                                                     | :white_check_mark:       |           |
| `port`           | Port specifies the port to connect to the SMTP server.                                            |       | 25        |
| `username`       | Username specifies the username to use when connecting to the SMTP server.                        |       |           |
| `password`       | Password specifies the password to use when connecting to the SMTP server.                        |       |           |
| `auth`           | SMTP authentication method. Options: `Unknown`, `Plain`, `Login`, `CRAMMD5`, `None`, `OAuth2`.                |       | Unknown   |
| `fromAddress`    | From specifies the email address to use as the sender.                                            | :white_check_mark:       |           |
| `fromName`       | FromName specifies the name to use as the sender.                                                 |       |           |
| `toAddress`      | List of recipient e-mails.                                                                        | :white_check_mark:       |           |
| `clientHost`     | The client host name sent to the SMTP server during HELLO phase. If set to "auto", uses OS hostname. |       | auto      |
| `encryption`     | Encryption method. Options: `Auto`, `None`, `ExplicitTLS`, `ImplicitTLS`.                                 |       | Auto      |
| `useHTML`        | Whether the message being sent is in HTML format.                                                 |       | false     |
| `useStartTLS`    | Whether to use the STARTTLS command (if the server supports it).                                  |       | true      |
| `templateSubject`| The subject template for the email.                                                               |       |           |
| `templateBody`   | The body template for the email.                                                                  |       |  [Default message](getting-start.md#template-body-alert-message)         |
