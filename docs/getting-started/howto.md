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
  image: {{dockerImages.whoami}}
  baseTag: v1.9.0
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

In this example the image `{{dockerImages.whoami}}` will be updated every 12 hours with the latest minor version.

3 - Check the Image TAG:

```bash
kubectl get image demo'

NAME   IMAGE                  TAG
demo   {{dockerImages.whoami}}
```

But you can force the update by running the following command:

```bash
kubectl annotate image demo kimup.cloudavenue.io/action=refresh
```

The Image TAG is now updated:

```bash
NAME   IMAGE                  TAG
demo   {{dockerImages.whoami}}   v1.10.0
```

4 - Make a deployment with the image:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: whoami
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: whoami
  template:
    metadata:
      annotations:
        kimup.cloudavenue.io/enabled: "true"
      labels:
        app: whoami
    spec:
      containers:
        - name: whoami
          image: {{dockerImages.whoami}}
```

5 - Apply the deployment:

```bash
kubectl apply -f deployment.yaml
```

Now the deployment is running with the image `{{dockerImages.whoami}}:v1.10.0` define by your rules in the CRD `Image`.
