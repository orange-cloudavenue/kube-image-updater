---
hide:
  - toc
---

# Overview

!!! warning  "Project is in early development."

    This project is in early development and is not yet ready for production use.
    You are welcome to try it out and provide feedback, but be aware that the
    API may change at any time.

**kube-image-updater** (A.K.A. **kimup**, which is pronounced /kim up/) is a tool that helps you to update the image of a Kubernetes Deployment, StatefulSet, DaemonSet, or CronJob. It can be used to update the image of a single resource or multiple resources at once.

**kimup** is designed to be simple to use and easy to deploy. It is an kubernetes operator with custom resource definition (CRD) that allows you to define the image update strategy and schedule.

The project is composed of 3 main components:
**kimup-operator:**
:     The main component that reconcile `Image` and CRD definition.

**kimup:**
:     The component that updates TAG of the `Image` resource.

**kimup-webhook:**
:     The component that patch resource creation with rules defined in `Image` resource.

Basic example of usage:

```yaml
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Image
metadata:
  labels:
    app.kubernetes.io/name: kube-image-updater
    app.kubernetes.io/managed-by: kustomize
  name: demo
  namespace: default
spec:
  image: ghcr.io/orange-cloudavenue/kube-image-updater
  baseTag: v0.0.19
  triggers:
    - type: crontab
      value: "00 00 */12 * * *"
  rules:
    - name: Automatic update semver minor
      type: semver-minor
      actions:
        - type: apply
```

The `Image` resource defines the image to update, the base tag, the triggers, and the rules. In this example, the image `ghcr.io/orange-cloudavenue/kube-image-updater` will be updated every 12 hours with the latest minor version.

It is structured around the following concepts:

- **Triggers**: define when the image should be updated. (Multiple triggers can be defined)
- **Rules**: define how the image should be updated. (Multiple rules can be defined)
- **Actions**: define what should be done after the image is updated (rule matched). (Multiple actions can be defined)

## Features

- **baseTag**: version of the image to update and on which we will apply the rules.

Type of **Triggers**:

- **crontab**: Update the image based on a crontab schedule.

Type of **Rules**:

- **semver-major**: Update the major version of the image.
- **semver-minor**: Update the minor version of the image.
- **semver-patch**: Update the patch version of the image.
- **regex**: Update the image based on a regex pattern.

Type of **Actions**:

- **refresh**: Apply the new image to the resource.
- **notify**: Notify a webhook with the new image.(Not implemented yet)
- **request-approval**: Request approval to apply the new image.(Not implemented yet)
<<<<<<< HEAD
=======

>>>>>>> e2a5d09 (chore: Add Doc (#17))
