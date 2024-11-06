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

## Configuration

Kimup Operator uses a dedicated kimup CRD to create and manage image resources. The CRD allows various configurations to define the behaviour of the image. See [docs.crds.dev](https://doc.crds.dev/github.com/orange-cloudavenue/kube-image-updater/kimup.cloudavenue.io/Image/v1alpha1) for more information about the Image CRD.

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
## Status

The following status can be set on an image:

| Last-Sync state | Description |
| --------------- | ----------- |
| &#34;ActionError&#34; | Status of the image when it is last sync error action. |
| &#34;Error&#34; | Status of the image when an error occurred. |
| &#34;GetImageError&#34; | Status of the image when it is last sync get error. |
| &#34;GetRuleError&#34; | Status of the image when it is last sync error get rule. |
| &#34;PullSecretsError&#34; | Status of the image when it is last sync error secrets. |
| &#34;RegistryError&#34; | Status of the image when it is last sync error registry. |
| &#34;Scheduled&#34; | Status of the image when it is last sync is scheduled. |
| &#34;Success&#34; | Status of the image when it is last sync success. |
| &#34;TagsError&#34; | Status of the image when it is last sync error tags. |


