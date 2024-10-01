---
hide:
  - toc
---

# HowTo

## How to Use

1 - Create an `Image` resource:

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

2 - Apply the `Image` resource:

```bash
kubectl apply -f image.yaml
```

In this example the image `ghcr.io/orange-cloudavenue/kube-image-updater` will be updated every 12 hours with the latest minor version.

3 - Check the Image TAG:

```bash
kubectl get image demo'

NAME   IMAGE                  TAG
demo   ghcr.io/azrod/golink
```

But you can force the update by running the following command:

```bash
kubectl annotate image demo kimup.cloudavenue.io/action=refresh
```

The Image TAG is now updated:

```bash
NAME   IMAGE                  TAG
demo   ghcr.io/azrod/golink   v0.1.0
```

4 - Make a deployment with the image:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: golink
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: golink
  template:
    metadata:
      annotations:
        kimup.cloudavenue.io/enabled: "true"
      labels:
        app: golink
    spec:
      containers:
        - name: golink
          image: ghcr.io/azrod/golink
          ports:
            - containerPort: 8080
```

5 - Apply the deployment:

```bash
kubectl apply -f deployment.yaml
```

Now the deployment is running with the image `ghcr.io/azrod/golink:v0.1.0` define by your rules in the CRD `Image`.
