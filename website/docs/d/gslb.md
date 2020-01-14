---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_gslb"
subcategory: "Global"
description: |-
  Get information about an existing GSLB.
---

# Data Source: sakuracloud_gslb

Get information about an existing GSLB.

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

* `id` - The id of the GSLB.
* `description` - The description of the GSLB.
* `fqdn` - The FQDN for accessing to the GSLB. This is typically used as value of CNAME record.
* `health_check` - A list of `health_check` blocks as defined below.
* `icon_id` - The icon id attached to the GSLB.
* `name` - The name of the GSLB.
* `server` - A list of `server` blocks as defined below.
* `sorry_server` - The IP address of the SorryServer. This will be used when all servers are down.
* `tags` - Any tags assigned to the GSLB.
* `weighted` - The flag to enable weighted load-balancing.


---

A `health_check` block exports the following:

* `delay_loop` - The interval in seconds between checks.
* `host_header` - The value of host header send when checking by HTTP/HTTPS.
* `path` - The path used when checking by HTTP/HTTPS.
* `port` - The port number used when checking by TCP.
* `protocol` - The protocol used for health checks. This will be one of [`http`/`https`/`tcp`/`ping`].
* `status` - The response-code to expect when checking by HTTP/HTTPS.

---

A `server` block exports the following:

* `enabled` - The flag to enable as destination of load balancing.
* `ip_address` - The IP address of the server.
* `weight` - The weight used when weighted load balancing is enabled.



