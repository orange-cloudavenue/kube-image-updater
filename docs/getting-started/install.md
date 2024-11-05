---
hide:
  - toc
---

# Install Kimup Operator

Kimup Operator is a Kubernetes operator used to manage images and their lifecycle, manage kimup-controller deployments. The operator is required for the functioning of the Kimup.

Resources managed by Kimup Operator are:

* [Image](../crd/image.md)
* [AlertConfig](../crd/alertconfig.md)
* [Kimup](../crd/kimup.md)

## Prerequisites

* A Kubernetes cluster with a version >= 1.28
* `kubectl` with kustomize installed and configured to connect to your cluster
* `cert-manager` installed in your cluster (See [cert-manager documentation](https://cert-manager.io/docs/installation/kubernetes/))

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

!!! warning "Namespace"
    For the moment only the `kimup-operator` namespace is supported.

### Deploy `kimup-controller`

For deploying `kimup-controller`, create a `Kimup` resource in the `kimup-operator` namespace:

```yaml
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Kimup
metadata:
  labels:
    app.kubernetes.io/name: kube-image-updater
  name: kimup
  namespace: kimup-operator
spec:
  name: demo
  logLevel: info
```

```bash
kubectl apply -f kimup.yaml
```

```bash
kubectl get kimup

NAME    STATE
kimup   ready
```