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
| kimup_actions_executed_duration                | The duration in seconds of action performed.                |
| kimup_actions_executed_error_total             | The total number of action performed with error.            |
| kimup_actions_executed_total                   | The total number of action performed.                       |
| kimup_admission_controller_patch_duration      | The duration in seconds of patch in admission controller.   |
| kimup_admission_controller_patch_error_total   | The total number of patch action performed with error.      |
| kimup_admission_controller_patch_total         | The total number of patch action performed.                 |
| kimup_admission_controller_request_duration    | The duration in seconds of request in admission controller. |
| kimup_admission_controller_request_error_total | The total number of request received with error.            |
| kimup_admission_controller_request_total       | The total number of request received.                       |
| kimup_events_triggerd_error_total              | The total number of events triggered with error.            |
| kimup_events_triggered_duration                | The duration in seconds of events triggered.                |
| kimup_events_triggered_total                   | The total number of events triggered.                       |
| kimup_registry_request_duration                |                                                             |
| kimup_registry_request_error_total             | The total number of registry evaluated with error.          |
| kimup_registry_request_total                   | The total number of registry evaluated.                     |
| kimup_rules_evaluated_duration                 | The duration in seconds of rules evaluated.                 |
| kimup_rules_evaluated_error_total              | The total number of rules evaluated with error.             |
| kimup_rules_evaluated_total                    | The total number of rules evaluated.                        |
| kimup_tags_available_sum                       | The total number of tags available for an image.            |
| kimup_tags_request_duration                    | The duration in seconds of the request to list tags.        |
| kimup_tags_request_error_total                 | The total number returned an error when calling list tags.  |
| kimup_tags_request_total                       | The total number of requests to list tags.                  |

