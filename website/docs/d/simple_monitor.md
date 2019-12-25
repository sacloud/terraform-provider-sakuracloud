---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_monitor"
subcategory: "Global"
description: |-
  Get information about an existing SimpleMonitor.
---

# Data Source: sakuracloud_simple_monitor

Get information about an existing SimpleMonitor.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more `condition` blocks as defined below.
* `id` - (Optional) .
* `names` - (Optional) .
* `tags` - (Optional) .

---

A `condition` block supports the following:

* `name` - (Required) .
* `values` - (Required) .

---

A `health_check` block supports the following:

* `password` - (Optional) .
* `username` - (Optional) .


## Attribute Reference

* `id` - The ID of the SimpleMonitor.
* `delay_loop` - .
* `description` - .
* `enabled` - .
* `health_check` - A list of `health_check` blocks as defined below.
* `icon_id` - .
* `notify_email_enabled` - .
* `notify_email_html` - .
* `notify_slack_enabled` - .
* `notify_slack_webhook` - .
* `tags` - .
* `target` - .


---

A `health_check` block exports the following:

* `community` - .
* `excepcted_data` - .
* `host_header` - .
* `oid` - .
* `path` - .
* `port` - .
* `protocol` - .
* `qname` - .
* `remaining_days` - .
* `sni` - .
* `snmp_version` - .
* `status` - .



