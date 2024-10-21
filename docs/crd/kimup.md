---
hide:
  - toc
---

# Custom Resource Definition `Kimup`

This is a custom resource definition for a Kimup. It is used to manage a deployment of a kimup-controller.

## Basic example

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

```

## Configuration

Kimup Operator uses a dedicated kimup CRD to create and manage kimup resources. The CRD allows various configurations to define the behaviour of the kimup controller. See [docs.crds.dev](https://doc.crds.dev/github.com/orange-cloudavenue/kube-image-updater/kimup.cloudavenue.io/Kimup/v1alpha1) for more information about the Kimup CRD.
