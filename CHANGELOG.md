## 0.1.0 (Unreleased)
### :rotating_light: **Breaking Changes**

* `crd kimup` - Refactor the `kimup` CRD to remove the `admissioncontroller` field. (GH-85)
* `kimup-admission-controller` - Remove component admission controller. This component is moved to `kimup-operator`, the `kimup-admission-controller` is no longer available. (GH-59)

### :rocket: **New Features**

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
