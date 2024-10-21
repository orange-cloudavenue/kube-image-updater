---
hide:
  - toc
---

# Metrics

kimup exposes metrics to monitor the performance. The metrics are exposed in the Prometheus format and can be scraped by Prometheus or any other monitoring tool that can scrape Prometheus.

## Settings

The following arguments can be used to configure the metrics *(Available in kimup-operator, kimup-controller and kimup-admission-controller)*:

| Flag           | Default  | Description               |
| -------------- | -------- | ------------------------- |
| --metrics      | false    | Enable metrics collection |
| --metrics-port | :9080    | Port to expose metrics on |
| --metrics-path | /metrics | Path to expose metrics on |


## Metrics

The following metrics are exposed:

| Metrics                                        | Description                                                 |
| ---------------------------------------------- | ----------------------------------------------------------- |
| kimup_actions_duration                         | The duration in seconds of action performed.                |
| kimup_actions_err_total                        | The total number of action performed with error.            |
| kimup_actions_total                            | The total number of action performed.                       |
| kimup_admission_controller_patch_duration      | The duration in seconds of patch in admission controller.   |
| kimup_admission_controller_patch_err_total     | The total number of patch action performed with error.      |
| kimup_admission_controller_patch_total         | The total number of patch action performed.                 |
| kimup_admission_controller_request_duration    | The duration in seconds of request in admission controller. |
| kimup_admission_controller_request_error_total | The total number of request received with error.            |
| kimup_admission_controller_request_total       | The total number of request received.                       |
| kimup_events_duration                          | The duration in seconds of events performed.                |
| kimup_events_err_total                         | The total number of events performed with error.            |
| kimup_events_total                             | The total number of events performed.                       |
| kimup_registry_duration                        | The duration in seconds of registry evaluated.              |
| kimup_registry_err_total                       | The total number of registry evaluated with error.          |
| kimup_registry_total                           | The total number of registry evaluated.                     |
| kimup_rules_duration                           | The duration in seconds of rules evaluated.                 |
| kimup_rules_err_total                          | The total number of rules evaluated with error.             |
| kimup_rules_total                              | The total number of rules evaluated.                        |
| kimup_tags_duration                            | The duration in seconds for func tags to list the tags.     |
| kimup_tags_total                               | The total number of func tags is called to list tags.       |
| kimup_tags_total_err                           | The total number return by the func tags with error.        |

