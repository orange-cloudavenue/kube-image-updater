
# Registry

## Generate htpasswd secret

```sh
docker run --rm --entrypoint htpasswd registry:2.7.0 -Bbn myuser mypasswd > htpasswd
```

Output:

```txt
myuser:$2y$05$P7carSLFLt2fvJzZbj10yuD0zc69z/dD9X92rzlbdKp0cr8Sjq9Km
```

## Create secret for registry

```sh
kubectl --namespace tests create secret generic auth-secret --from-file=htpasswd --dry-run=client -o yaml > secret.yaml
```
