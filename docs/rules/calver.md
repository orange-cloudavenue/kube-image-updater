---
hide:
  - toc
---

#  Calendar Versioning

The `calver` rule allows you to define a rule that will be executed when the image is updated with a new calver version.
It follows the [Calendar Versioning](https://calver.org/) specification.
A lot of options are available to match the version you want to update.
The format of the version [following this logic regex](https://regex101.com/r/25eVYJ/2).

* `calver-major`: Update the image with the latest major version.
* `calver-minor`: Update the image with the latest minor version.
* `calver-patch`: Update the image with the latest patch version.

**`calver-major`** is the most restrictive and will only update the image when the major version is updated [calver documentation](https://calver.org/).
``` { .yaml .no-copy title="calver rule" }
version: 2024.0.0
Match: >=2024.*.* # (1)
```

1.  :man_raising_hand: For more information about the calver range, you can check the [calver documentation](https://calver.org/).

**`calver-minor`** is less restrictive and will update the image when the minor version is updated.
``` { .yaml .no-copy title="calver rule" }
version: 2024.0.0
Match: >=2024.1.* <2 # (1)
```

1. :man_raising_hand: For more information about the calver range, you can check the [calver documentation](https://calver.org/).

**`calver-patch`** is the least restrictive and will update the image when the patch version is updated.
``` { .yaml .no-copy title="calver rule" }
version: 2024.0.0
Match: >=2024.0.1 <2024.1.0 # (1)
```

1. :man_raising_hand: For more information about the calver range, you can check the [calver documentation](https://calver.org/).

## Who to use

Create an `Image` resource with the `calver` rule.

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
  baseTag: v2024.0.4
  triggers:
    - [...]
  rules:
    - name: Notify when calver major is detected
      type: calver-major
      actions:
        - type: alert-xxx
          [...]
    - name: Automatic update calver minor
      type: calver-minor
      actions:
        - type: apply
    - name: Automatic update calver patch
      type: calver-patch
      actions:
        - type: apply
```
