---
hide:
  - toc
---

# Scope

## Overview

Kimup operator allow you to manage at what level it monitors to apply or not the tag when creating a pod.

Scope is defined by the annotation `kimup.cloudavenue.io/enabled` on the namespace or the pod.

## Logical

![Logical pod creation schema](logical-pod-creation-light.png#only-light)
![Logical pod creation schema](logical-pod-creation-dark.png#only-dark)

## Namespace

When the annotation `kimup.cloudavenue.io/enabled: "true"` is set on a namespace, the operator will only apply the tag on the pods that are created in this namespace.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "true"
```

## Pod

When the annotation `kimup.cloudavenue.io/enabled: "true"` is set on a pod, the operator will only apply the tag on this pod.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: your-pod
  namespace: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "true"
```

## Ignore for a pod

When the annotation `kimup.cloudavenue.io/enabled: "false"` is set on a pod, the operator will ignore this pod even if the namespace is enabled.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: your-pod
  namespace: your-env
  annotations:
    kimup.cloudavenue.io/enabled: "false"
```
