---
hide:
  - toc
---

# Always

!!! warning "Only for testing purpose"
    The `always` trigger is only for testing purpose and should not be used in production.

The `always` rule allows you to define a rule that will be executed every time the refresh is triggered.

## Who to use

Create an `Image` resource with the `always` rule.

```yaml hl_lines="13-15"
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
    - name: Always update
      type: always
      actions:
        - [...]
```
