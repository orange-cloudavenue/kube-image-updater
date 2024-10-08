---
hide:
  - toc
---

# Regex

The `regex` rule allows you to define a rule that will be executed when the image is updated with a new version that matches a regular expression.
Regex is compatible with [Golang Regex format](https://pkg.go.dev/regexp/syntax). Use the great [regex101.com](https://regex101.com/) website to test your regular expression.

## Who to use

Create an `Image` resource with the `regex` rule.

```yaml hl_lines="14-18"
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Image
metadata:
  labels:
    app.kubernetes.io/name: kube-image-updater
    app.kubernetes.io/managed-by: kustomize
  name: image-sample-with-auth
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

1. For more information about this regular expression, you can check the [regex101.com/r/prt9tw/1](https://regex101.com/r/prt9tw/1).
