---
hide:
  - toc
---

#  Calendar Versioning

The `calver` rule allows you to define a rule that will be executed when the image is updated with a new calver version.
It follows the [Calendar Versioning](https://calver.org/) specification.
A lot of options are available to match the version you want to update.

The format of the version [following this logic regex](https://regex101.com/r/dRh6UI/1).
```regex
^([0-9]{4}|[0-9]{2})(\.[0-9]{1,2})?(\.[0-9]{1,2})?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$
```

Format allowed:
```s
# YYYY.MM.XX
2024
2024.01
2024.01.01

# YY.MM.XX
24
24.01

# YY.M.X
24.1
24.1.1

2024.01.01-dev.1
```


* `calver-major`: Update the image with the latest major version.
* `calver-minor`: Update the image with the latest minor version.
* `calver-patch`: Update the image with the latest patch version.
* `calver-prerelease`: Update the image with the latest prerelease version.

**`calver-major`** is the most restrictive and will only update the image when the major version is updated [calver documentation](https://calver.org/).
Most of time the major is a year date (eg: `2024`, `2025`, `2026`, ...).
``` { .yaml .no-copy title="calver rule" }
version: 2024.0.0
Match: >=2025.*.* # (1)
```

1.  :man_raising_hand: For more information about the calver range, you can check the [calver documentation](https://calver.org/).

**`calver-minor`** is less restrictive and will update the image when the minor version is updated.
Most of time the minor is a month date (eg: `2024.01`, `2024.02`, `2024.03`, ...).
``` { .yaml .no-copy title="calver rule" }
version: 2024.0.0
Match: >=2024.1.* and <2025.0.0 # (1)
```

1. :man_raising_hand: For more information about the calver range, you can check the [calver documentation](https://calver.org/).

**`calver-patch`** is the least restrictive and will update the image when the patch version is updated.
Most of time the patch is a day date (eg: `2024.01.01`, `2024.01.02`, `2024.01.03`, ...).
``` { .yaml .no-copy title="calver rule" }
version: 2024.0.0
Match: >=2024.0.1 and < 2024.1.0 # (1)
```

**`calver-prerelease`** is the least restrictive and will update the image when the prerelease version is updated.
The prerelease version is the part after the `-xxx.` in the version. It's an incremental number. (eg: `2024.0.0-dev.1`, `2024.0.0-dev.2`, `2024.0.0-dev.3`,...).
``` { .yaml .no-copy title="calver rule" }
version: 2024.0.0-dev.0
Match: >=2024.0.0-dev.1
```

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
