---
hide:
  - toc
---

# Apply

The `apply` action update the new image tag to the resource.

## Who to use

Create an `Image` resource with the `apply` action.

```yaml hl_lines="18"
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
      value: "^v?[0-9].[0-9].[0-9]-dev[0-9]$" # (1)
      actions:
        - type: apply
```

In this example, the `apply` action will be executed when the image is updated with a new version that matches the regular expression `^v?[0-9].[0-9].[0-9]-dev[0-9]$`.
