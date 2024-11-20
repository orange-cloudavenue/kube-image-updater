## 0.3.0 (Unreleased)

### :bug: **Bug Fixes**

* `kimup-controller` - Fix rbac authorization to read secrets. (GH-121)

### :dependabot: **Dependencies**

* deps: bumps golang.org/x/term from 0.25.0 to 0.26.0 (GH-124)

## 0.2.0 (November  7, 2024)

### :rocket: **New Features**

* `Image` - Add new annotation `kimup.cloudavenue.io/failure-policy` to control the failure policy of the image tag mutation. The default value is `Fail`. The supported values are `Fail` and `Ignore`. (GH-93)

### :tada: **Improvements**

* `image/status` - Fix the documentation of the `image/status. (GH-101)

### :bug: **Bug Fixes**

* `healthz` - Fix the issue where the port specified in the `healthz-port` configuration was not being used. (GH-100)
* `image` - Now the image is displayed correctly the `last-result`. (GH-101)
* `image` - Now the status of the image is displayed correctly in the `image` list if the image/repository is not found in the registry. (GH-103)
* `kimup-controller` - Now the application restart if the kubernetes watch connection is lost. (GH-109)
* `kimup-operator` - Now webhook start with default port 9443. (GH-105)
* `metrics` - Fix the issue where the port specified in the `metrics-port` configuration was not being used. (GH-100)

### :dependabot: **Dependencies**

* deps: bumps github.com/onsi/ginkgo/v2 from 2.20.2 to 2.21.0 (GH-96)

## 0.1.0 (November  3, 2024)
### :rotating_light: **Breaking Changes**

* `crd kimup` - Refactor the `kimup` CRD to remove the `admissioncontroller` field. (GH-85)
* `kimup-admission-controller` - Remove component admission controller. This component is moved to `kimup-operator`, the `kimup-admission-controller` is no longer available. (GH-59)

### :rocket: **New Features**

* `chore` - Now the mutating webhook configuration used for mutate image tag on pod creation is created by the operator itself. This is done to avoid the need for the user to create the mutating webhook configuration manually. The operator will also update the mutating webhook configuration if the user changes the configuration (annotations) in the namespace. (GH-87)
* `feat` - Add calver Rule Semantic. (GH-58)
* `kimup-operator` - Now the `kimup-operator` allow to set scope of the image-tag mutation. The scope can be set to `Namespaced` or `Pod`. (GH-59)

### :tada: **Improvements**

* `operator` - Now operator expose /healthz for health check and /readyz for readiness check. (GH-89)
* `operator` - Now operator expose `metrics` for monitoring. (GH-89)
* `webserver` - Now all webserver use same `log` format. (GH-89)

### :dependabot: **Dependencies**

* deps: bumps crazy-max/ghaction-setup-docker from 3.3.0 to 3.4.0 (GH-83)
* deps: bumps github.com/opencontainers/runc from 1.1.13 to 1.1.14 (GH-65)
* deps: bumps github.com/prometheus/client_golang from 1.20.4 to 1.20.5 (GH-73)
* deps: bumps golang.org/x/term from 0.23.0 to 0.25.0 (GH-74)
* deps: bumps k8s.io/client-go from 0.31.1 to 0.31.2 (GH-79)
* deps: bumps sigs.k8s.io/controller-runtime from 0.19.0 to 0.19.1 (GH-82)

## 0.0.5 (October 18, 2024)
## 0.0.2 (October 18, 2024)

## 0.0.1 (October 17, 2024)

### :rocket: Initial release
