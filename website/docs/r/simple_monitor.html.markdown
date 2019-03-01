---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_monitor"
sidebar_current: "docs-sakuracloud-resource-global-simple-monitor"
description: |-
  Provides a SakuraCloud Simple Monitor resource. This can be used to create, update, and delete Simple Monitors.
---

# sakuracloud\_simple\_monitor

Provides a SakuraCloud Simple Monitor resource. This can be used to create, update, and delete Simple Monitors.

## Example Usage

```hcl
# Create a new Simple Monitor(protocol: https)
resource "sakuracloud_simple_monitor" "foobar" {
  target = "www.example.com"
  health_check {
    protocol    = "https"
    delay_loop  = 60
    path        = "/"
    status      = "200"
    host_header = "hostname.example.com"
    sni         = true
  }

  # for Basic Auth
  # username = "foo"
  # password = "bar"

  notify_email_enabled = true
  notify_email_html    = true
  notify_slack_enabled = true
  notify_slack_webhook = "https://hooks.slack.com/services/XXX/XXX/XXXXXX"

  description = "description"
  tags        = ["foo", "bar"]
}

# Create a new Simple Monitor(protocol: sslcertificate)
resource "sakuracloud_simple_monitor" "cert" {
  target = "www.example.com"
  health_check {
    protocol       = "sslcertificate"
    remaining_days = 30
  }
}

```

## Argument Reference

The following arguments are supported:

* `target` - (Required) The HostName or IP address of monitoring target.
* `health_check` - (Required) Health check rules. It contains some attributes to [Health Check](#health-check).
* `notify_email_enabled` - (Optional) The flag of enable/disable notification by E-mail.
* `notify_email_html` - (Optional) The flag of enable/disable HTML format for E-mail.
* `notify_slack_enabled` - (Optional) The flag of enable/disable notification by slack.
* `notify_slack_webhook` - (Optional) The webhook URL of destination of slack notification.
* `enabled` - (Optional) The flag of enable/disable monitoring.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon of the resource.

### Health Check

Attributes for Health Check:

* `protocol` - (Required) Protocol used in health check.  
Valid value is one of the following: [ "http" / "https" / "ping" / "tcp" / "dns" / "ssh" / "smtp" / "pop3" / "snmp" / "sslcertificate" ]
* `delay_loop` - (Optional) Health check access interval (unit:`second`). 
* `host_header` - (Optional) The value of `Host` header used in http/https health check access.
* `path` - (Optional) The request path used in http/https health check access.
* `status` - (Optional) HTTP status code expected by health check access.
* `sni` - (Optional) The flag of enable/disable SNI.
* `username` - (Optional) The Basic Auth Username used in http/https health check access.
* `password` - (Optional) The Basic Auth Password request path used in http/https health check access.
* `port` - (Optional) Port number used in health check access.
* `qname` - (Optional) The QName value used in dns health check access.
* `excepcted_data` - (Optional) The expect value used in dns/snmp health check.
* `community` - (Optional) The community name used in snmp health check.
* `snmp_version` - (Optional) SNMP cersion used in snmp health check.
* `oid` - (Optional) The OID used in snmp health check.
* `remaining_days` - (Optional) The number of remaining days used in ssh-certificate check.

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

## Import

Simple Monitors can be imported using the Simple Monitor ID.

```
$ terraform import sakuracloud_simple_monitor.foobar <simple_monitor_id>
```
