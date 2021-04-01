---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_simple_monitor"
subcategory: "Global"
description: |-
  Get information about an existing Simple Monitor.
---

# Data Source: sakuracloud_simple_monitor

Get information about an existing Simple Monitor.

## Example Usage

```hcl
data "sakuracloud_simple_monitor" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.


## Attribute Reference

* `id` - The id of the Simple Monitor.
* `delay_loop` - The interval in seconds between checks.
* `description` - The description of the SimpleMonitor.
* `enabled` - The flag to enable monitoring by the simple monitor.
* `health_check` - A list of `health_check` blocks as defined below.
* `icon_id` - The icon id attached to the SimpleMonitor.
* `notify_email_enabled` - The flag to enable notification by email.
* `notify_email_html` - The flag to enable HTML format instead of text format.
* `notify_interval` - The interval in hours between notification.
* `notify_slack_enabled` - The flag to enable notification by slack/discord.
* `notify_slack_webhook` - The webhook URL for sending notification by slack/discord.
* `tags` - Any tags assigned to the SimpleMonitor.
* `target` - The monitoring target of the simple monitor. This will be IP address or FQDN.


---

A `health_check` block exports the following:

* `community` - The SNMP community string used when checking by SNMP.
* `excepcted_data` - The expected value used when checking by DNS.
* `host_header` - The value of host header send when checking by HTTP/HTTPS.
* `http2` - The flag to enable HTTP/2 when checking by HTTPS.
* `oid` - The SNMP OID used when checking by SNMP.
* `password` - The password for basic auth used when checking by HTTP/HTTPS.
* `path` - The path used when checking by HTTP/HTTPS.
* `port` - The target port number.
* `protocol` - The protocol used for health checks. This will be one of [`http`/`https`/`ping`/`tcp`/`dns`/`ssh`/`smtp`/`pop3`/`snmp`/`sslcertificate`].
* `qname` - The FQDN used when checking by DNS.
* `remaining_days` - The number of remaining days until certificate expiration used when checking SSL certificates.
* `sni` - The flag to enable SNI when checking by HTTP/HTTPS.
* `snmp_version` - The SNMP version used when checking by SNMP.
* `status` - The response-code to expect when checking by HTTP/HTTPS.
* `username` - The user name for basic auth used when checking by HTTP/HTTPS.


