---
hide:
  - toc
---

# Failure Policy

## Overview

Kimup operator allows you to manage the behavior of the operator when it fails to apply the tag on a pod. The failure policy is defined by the annotation `kimup.cloudavenue.io/failure-policy` on the **namespace** or the **pod**. The failure policy can be set to `fail` or `ignore`. The default value is `fail`.

!!! warning
    The annotation `kimup.cloudavenue.io/enabled` must be set to `true` on the namespace or the pod to apply the failure policy. If the annotation is not set, the failure policy will be ignored. See [Scope](../getting-started/scope.md) for more information.

## Logical

![Logical pod creation schema](../getting-started/logical-pod-creation-light.png#only-light)
![Logical pod creation schema](../getting-started/logical-pod-creation-dark.png#only-dark)

## Apply the failure policy

When the annotation `kimup.cloudavenue.io/failure-policy: "fail"` is set on a namespace, the operator will fail if it can't apply the tag on a pod created in this namespace.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "true"
    kimup.cloudavenue.io/failure-policy: "fail"
```

When the annotation `kimup.cloudavenue.io/failure-policy: "ignore"` is set on a namespace, the operator will ignore the failure if it can't apply the tag on a pod created in this namespace.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "true"
    kimup.cloudavenue.io/failure-policy: "ignore"
```

For a pod, the same logic applies.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: your-pod
  namespace: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "true"
    kimup.cloudavenue.io/failure-policy: "fail"
```

## Override the failure policy for a pod

When the annotation `kimup.cloudavenue.io/failure-policy` is set on a namespace, the operator will apply the failure policy on all pods created in this namespace. If the annotation is set on a pod, the operator will apply the failure policy defined on the pod, but i will be necessary to set the annotation `kimup.cloudavenue.io/enabled` to `true` on the pod.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "true"
    kimup.cloudavenue.io/failure-policy: "fail"
---
apiVersion: v1
kind: Pod
metadata:
  name: your-pod
  namespace: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "true"
    kimup.cloudavenue.io/failure-policy: "ignore"
```
