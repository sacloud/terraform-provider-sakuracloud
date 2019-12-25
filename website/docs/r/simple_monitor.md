---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_monitor"
subcategory: "Global"
description: |-
  Manages a SakuraCloud SimpleMonitor.
---

# sakuracloud_simple_monitor

Manages a SakuraCloud SimpleMonitor.

## Argument Reference

* `delay_loop` - (Optional) . Defaults to `60`.
* `description` - (Optional) .
* `enabled` - (Optional) . Defaults to `true`.
* `health_check` - (Required) A `health_check` block as defined below.
* `icon_id` - (Optional) .
* `notify_email_enabled` - (Optional) . Defaults to `true`.
* `notify_email_html` - (Optional) .
* `notify_interval` - (Optional) Unit: Hours. Defaults to `2`.
* `notify_slack_enabled` - (Optional) .
* `notify_slack_webhook` - (Optional) .
* `tags` - (Optional) .
* `target` - (Required) . Changing this forces a new resource to be created.


---

A `health_check` block supports the following:

* `community` - (Optional) .
* `excepcted_data` - (Optional) .
* `host_header` - (Optional) .
* `oid` - (Optional) .
* `password` - (Optional) .
* `path` - (Optional) .
* `port` - (Optional) .
* `protocol` - (Required) .
* `qname` - (Optional) .
* `remaining_days` - (Optional) .
* `sni` - (Optional) .
* `snmp_version` - (Optional) .
* `status` - (Optional) .
* `username` - (Optional) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the SimpleMonitor

* `read` -   (Defaults to 5 minutes) Used when reading the SimpleMonitor

* `update` - (Defaults to 5 minutes) Used when updating the SimpleMonitor

* `delete` - (Defaults to 5 minutes) Used when deregistering SimpleMonitor



## Attribute Reference

* `id` - The ID of the SimpleMonitor.




