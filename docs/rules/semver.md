---
hide:
  - toc
---

#  Semantic Versioning (semver)

The `semver` rule allows you to define a rule that will be executed when the image is updated with a new semver version.
It follows the [Semantic Versioning](https://semver.org/) specification.
A lot of options are available to match the version you want to update.

* `semver-major`: Update the image with the latest major version.
* `semver-minor`: Update the image with the latest minor version.
* `semver-patch`: Update the image with the latest patch version.

**`semver-major`** is the most restrictive and will only update the image when the major version is updated.
``` { .yaml .no-copy title="Semver rule" }
version: 1.0.0
Match: >=2.*.* # (1)
```

1.  :man_raising_hand: For more information about the semver range, you can check the [semver documentation](https://semver.org/#spec-item-8).

**`semver-minor`** is less restrictive and will update the image when the minor version is updated.
``` { .yaml .no-copy title="Semver rule" }
version: 1.0.0
Match: >=1.1.* <2 # (1)
```

1. :man_raising_hand: For more information about the semver range, you can check the [semver documentation](https://semver.org/#spec-item-7).

**`semver-patch`** is the least restrictive and will update the image when the patch version is updated.
``` { .yaml .no-copy title="Semver rule" }
version: 1.0.0
Match: >=1.0.1 <1.1.0 # (1)
```

1. :man_raising_hand: For more information about the semver range, you can check the [semver documentation](https://semver.org/#spec-item-6).

## Who to use

Create an `Image` resource with the `semver` rule.

```yaml hl_lines="15 16 17 20 21 22 24 25 26"
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
    - name: Notify when semver major is detected
      type: semver-major
      actions:
        - type: alert-xxx
          [...]
    - name: Automatic update semver minor
      type: semver-minor
      actions:
        - type: apply
    - name: Automatic update semver patch
      type: semver-patch
      actions:
        - type: apply
```
