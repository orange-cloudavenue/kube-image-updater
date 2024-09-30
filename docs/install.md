# Install

## How to Install (In Development...)

1 - **Install the CRD:**

Git clone project and run the following command:
```bash
kubectl apply -f config/crd/bases/kimup.cloudavenue.io_images.yaml
```

2 - Start all services:

For Operator
```bash
make run-operator
```

For Webhook
```bash
go run ./cmd/webhook --insideCluster=false
```

For kimup
```bash
go run ./cmd/kimup --insideCluster=false
```
