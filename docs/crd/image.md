---
hide:
  - toc
---

# Custom Resource Definition `Image`

This is a custom resource definition for an image. It is used to store information about an image.
`Image` is a namespaced resource.

## Basic example

```yaml
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Image
metadata:
  name: image-sample
spec:
  image: alpine
  baseTag: v1.0.0
  triggers:
    - <trigger>
    - <trigger>
  rules:
    - <rule>
    - <rule>
```

## Advanced

### Use authenticated registry

Use the `imagePullSecrets` field to specify the name of the secret to use to authenticate with the registry.

```yaml
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Image
metadata:
  name: image-sample
spec:
    image: custom-registry.io/image
    baseTag: v1.0.0
    imagePullSecrets:
        - name:  registry-local
    triggers:
        - <trigger>
        - <trigger>
    rules:
        - <rule>
        - <rule>
```

### Self-signed certificate

Use the `insecureSkipTLSVerify` field to skip the verification of the TLS certificate.

```yaml
kind: Image
metadata:
  name: image-sample
spec:
    image: custom-registry.io/image
    baseTag: v1.0.0
    insecureSkipTLSVerify: true
    triggers:
        - <trigger>
        - <trigger>
    rules:
        - <rule>
        - <rule>
```
