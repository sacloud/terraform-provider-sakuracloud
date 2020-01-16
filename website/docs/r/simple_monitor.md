---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_monitor"
subcategory: "Global"
description: |-
  Manages a SakuraCloud Simple Monitor.
---

# sakuracloud_simple_monitor

Manages a SakuraCloud Simple Monitor.

## Example Usage

```hcl
resource "sakuracloud_simple_monitor" "foobar" {
  target = "www.example.com"

  delay_loop = 60
  health_check {
    protocol    = "https"
    port        = 443
    path        = "/"
    status      = "200"
    host_header = "example.com"
    sni         = true
    # username    = "username"
    # password    = "password"
  }

  description = "description"
  tags        = ["tag1", "tag2"]

  notify_email_enabled = true
  notify_email_html    = true
  notify_slack_enabled = true
  notify_slack_webhook = "https://hooks.slack.com/services/xxx/xxx/xxx"
}
```

## Argument Reference

* `target` - (Required) The monitoring target of the simple monitor. This must be IP address or FQDN. Changing this forces a new resource to be created.
* `health_check` - (Required) A `health_check` block as defined below.
* `delay_loop` - (Optional) The interval in seconds between checks. This must be in the range [`60`-`3600`]. Default:`60`.
* `enabled` - (Optional) The flag to enable monitoring by the simple monitor. Default:`true`.

#### Notification

* `notify_email_enabled` - (Optional) The flag to enable notification by email. Default:`true`.
* `notify_email_html` - (Optional) The flag to enable HTML format instead of text format.
* `notify_interval` - (Optional) The interval in hours between notification. Default:`2`.
* `notify_slack_enabled` - (Optional) The flag to enable notification by slack/discord.
* `notify_slack_webhook` - (Optional) The webhook URL for sending notification by slack/discord.

#### Common Arguments

* `description` - (Optional) The description of the SimpleMonitor. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the SimpleMonitor.
* `tags` - (Optional) Any tags to assign to the SimpleMonitor.


---

A `health_check` block supports the following:

* `protocol` - (Required) The protocol used for health checks. This must be one of [`http`/`https`/`ping`/`tcp`/`dns`/`ssh`/`smtp`/`pop3`/`snmp`/`sslcertificate`].
* `port` - (Optional) The target port number.

#### DNS

* `excepcted_data` - (Optional) The expected value used when checking by DNS.
* `qname` - (Optional) The FQDN used when checking by DNS.

#### HTTP/HTTPS

* `host_header` - (Optional) The value of host header send when checking by HTTP/HTTPS.
* `password` - (Optional) The password for basic auth used when checking by HTTP/HTTPS.
* `username` - (Optional) The user name for basic auth used when checking by HTTP/HTTPS.
* `path` - (Optional) The path used when checking by HTTP/HTTPS.
* `sni` - (Optional) The flag to enable SNI when checking by HTTP/HTTPS.
* `status` - (Optional) The response-code to expect when checking by HTTP/HTTPS.

#### Certificate

* `remaining_days` - (Optional) The number of remaining days until certificate expiration used when checking SSL certificates. This must be in the range [`1`-`9999`].


#### SNMP 

* `community` - (Optional) The SNMP community string used when checking by SNMP.
* `oid` - (Optional) The SNMP OID used when checking by SNMP.
* `snmp_version` - (Optional) The SNMP version used when checking by SNMP. This must be one of `1`/`2c`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Simple Monitor
* `update` - (Defaults to 5 minutes) Used when updating the Simple Monitor
* `delete` - (Defaults to 5 minutes) Used when deleting Simple Monitor

## Attribute Reference

* `id` - The id of the Simple Monitor.

