<div align="center">
    <a href="https://github.com/orange-cloudavenue/kube-image-updater/releases/latest">
      <img alt="Latest release" src="https://img.shields.io/github/v/release/orange-cloudavenue/kube-image-updater?style=for-the-badge&logo=starship&color=C9CBFF&logoColor=D9E0EE&labelColor=302D41&include_prerelease&sort=semver" />
    </a>
    <a href="https://github.com/orange-cloudavenue/kube-image-updater/pulse">
      <img alt="Last commit" src="https://img.shields.io/github/last-commit/orange-cloudavenue/kube-image-updater?style=for-the-badge&logo=starship&color=8bd5ca&logoColor=D9E0EE&labelColor=302D41"/>
    </a>
    <a href="https://github.com/orange-cloudavenue/kube-image-updater/blob/main/LICENSE">
      <img alt="License" src="https://img.shields.io/github/license/orange-cloudavenue/kube-image-updater?style=for-the-badge&logo=starship&color=ee999f&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/orange-cloudavenue/kube-image-updater/stargazers">
      <img alt="Stars" src="https://img.shields.io/github/stars/orange-cloudavenue/kube-image-updater?style=for-the-badge&logo=starship&color=c69ff5&logoColor=D9E0EE&labelColor=302D41" />
    </a>
    <a href="https://github.com/orange-cloudavenue/kube-image-updater/issues">
      <img alt="Issues" src="https://img.shields.io/github/issues/orange-cloudavenue/kube-image-updater?style=for-the-badge&logo=bilibili&color=F5E0DC&logoColor=D9E0EE&labelColor=302D41" />
    </a>
</div>

# Kubernetes Image Updater

> [!CAUTION]
> This project is in early development and is not yet ready for production use. You are welcome to try it out and provide feedback, but be aware that theAPI may change at any time.

Kubernetes Image Updater is a kubernetes operator with custom resource definition (CRD) that allows you to define the image update strategy and schedule.

Useful links:

* [Kube Image Updater documentation](https://orange-cloudavenue.github.io/kube-image-updater/)

## Using the Kube Image Updater

To quickly get started with the Kube Image Updater, you can use the following example:

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

## Contributing

This project is open source and contributions are welcome.
You may also report issues or feature requests on the [GitHub issue tracker](https://github.com/orange-cloudavenue/kube-image-updater/issues/new/choose).
