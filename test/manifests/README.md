# Unit test env

## Setup

```bash
microk8s start
microk8s addons enable ingress
```

## Deploy

```bash
kustomize build --enable-helm . | kubectl apply -f -
```

## Remove

```bash
kustomize build --enable-helm . | kubectl delete -f -
```
