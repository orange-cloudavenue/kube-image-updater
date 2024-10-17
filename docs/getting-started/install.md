---
hide:
  - toc
---

# Install Kimup Operator

Kimup Operator is a Kubernetes operator used to manage images and their lifecycle, manage kimup-admission-controller and kimup-webhook deployments. The operator is required for the functioning of the Kimup.

Resources managed by Kimup Operator are:

* [Image](crd/image.md)
* [AlertConfig](crd/alertconfig.md)
* [Kimup](crd/kimup.md)

## Prerequisites

* A Kubernetes cluster with a version >= 1.19
* `kubectl` with kustomize installed and configured to connect to your cluster

## Installation

### Install custom resource definitions

```bash
kubectl apply -k "https://github.com/orange-cloudavenue/kube-image-updater/manifests/crd/?ref={{git.short_tag}}"
```

### Install Kimup Operator

```bash
kubectl apply -k "https://github.com/orange-cloudavenue/kube-image-updater/manifests/operator/?ref={{git.short_tag}}"
```

By default, Kimup Operator is installed in the `kimup-operator` namespace.

### Deploy `kimup-admission-controller` and `kimup-controller`

For deploying `kimup-admission-controller` and `kimup-controller`, create a `Kimup` resource:

```yaml
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Kimup
metadata:
  labels:
    app.kubernetes.io/name: kube-image-updater
  name: kimup
spec:
  controller:
    name: demo
    logLevel: info
  admissionController:
    name: demo
    logLevel: info
```
