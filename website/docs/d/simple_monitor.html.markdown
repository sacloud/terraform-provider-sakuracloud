---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_monitor"
sidebar_current: "docs-sakuracloud-datasource-simple-monitor"
description: |-
  Get information on a SakuraCloud Simple Monitor.
---

# sakuracloud\_simple\_monitor

Use this data source to retrieve information about a SakuraCloud Simple Monitor.

## Example Usage

```hcl
data "sakuracloud_simple_monitor" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

* `id` - The ID of the resource.
* `target` - The HostName or IP address of monitoring target.
* `health_check` - Health check rules. It contains some attributes to [Health Check](#health-check).
* `notify_email_enabled` - The flag of enable/disable notification by E-mail.
* `notify_email_html` - The flag of enable/disable HTML format for E-mail.
* `notify_slack_enabled` - The flag of enable/disable notification by slack.
* `notify_slack_webhook` - The webhook URL of destination of slack notification.
* `enabled` - The flag of enable/disable monitoring.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.

### Health Check

Attributes for Health Check:

* `protocol` - Protocol used in health check.
* `delay_loop` - Health check access interval (unit:`second`). 
* `host_header` - The value of `Host` header used in http/https health check access.
* `path` - The request path used in http/https health check access.
* `status` - HTTP status code expected by health check access.
* `sni` - The flag of enable/disable SNI.
* `username` - The Basic Auth Username used in http/https health check access.
* `password` - The Basic Auth Password request path used in http/https health check access.
* `port` - Port number used in health check access.
* `qname` - The QName value used in dns health check access.
* `excepcted_data` - The expect value used in dns/snmp health check.
* `community` - The community name used in snmp health check.
* `snmp_version` - SNMP cersion used in snmp health check.
* `oid` - The OID used in snmp health check.
* `remaining_days` - The number of remaining days used in ssh-certificate check.
