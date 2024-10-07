---
hide:
  - toc
---

# Crontab

The `crontab` trigger allows you to define a cron expression to schedule the image refresh.

The cron expression is a string representing a set of times, using 6 fields separated by white spaces. The fields represent:

1. Seconds (0-59)
2. Minutes (0-59)
3. Hours (0-23)
4. Day of month (1-31)
5. Month (1-12)
6. Day of week (0-6) (Sunday to Saturday)

`*` is a wildcard character that matches all values.

**Examples:**

* `00 00 */12 * * *` will trigger the image every 12 hours.
* `00 00 00 * * *` will trigger the image every day at midnight.
* `00 00 00 1 * *` will trigger the image every first day of the month at midnight.
* `00 00 00 * * 1` will trigger the image every Monday at midnight.


## Who to use

Create an `Image` resource with the `crontab` trigger.

Every 12 hours the image will execute rule defined in the `rules` section.

```yaml hl_lines="11-12"
apiVersion: kimup.cloudavenue.io/v1alpha1
kind: Image
metadata:
  labels:
    app.kubernetes.io/managed-by: kustomize
  name: demo
spec:
  image: registry.127.0.0.1.nip.io/demo
  baseTag: v0.0.4
  triggers:
    - type: crontab
      value: "00 00 */12 * * *"
  rules:
    - [...]
```
