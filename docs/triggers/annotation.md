---
hide:
  - toc
---

# Annotation

The `annotation` trigger allows you to define a Kubernetes annotation to schedule the image refresh.

## Who to use

```sh
kubectl annotate image demo kimup.cloudavenue.io/action=refresh
```

The `demo` image will be triggered instantly and launch the rule defined in the `rules` section.
