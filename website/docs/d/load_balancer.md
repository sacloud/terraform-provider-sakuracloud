---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
subcategory: "Appliance"
description: |-
  Get information about an existing Load Balancer.
---

# Data Source: sakuracloud_load_balancer

Get information about an existing Load Balancer.

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

* `id` - The id of the Load Balancer.
* `description` - The description of the load balancer.
* `gateway` - The IP address of the gateway used by load balancer.
* `high_availability` - The flag to enable HA mode.
* `icon_id` - The icon id attached to the load balancer.
* `ip_addresses` - The list of the IP address assigned to the load balancer.
* `name` - The name of the load balancer.
* `netmask` - The bit length of the subnet assigned to the load balancer.
* `plan` - The plan name of the load balancer. This will be one of [`standard`/`highspec`].
* `switch_id` - The id of the switch connected from the load balancer.
* `tags` - Any tags assigned to the load balancer.
* `vip` - A list of `vip` blocks as defined below.
* `vrid` - The Virtual Router Identifier. This is only used when `high_availability` is set `true`.
* `zone` - The name of zone that the load balancer is in (e.g. `is1a`,`tk1a`).


---

A `vip` block exports the following:

* `delay_loop` - The interval in seconds between checks.
* `description` - The description of the VIP.
* `port` - The target port number for load-balancing.
* `server` - A list of `server` blocks as defined below.
* `sorry_server` - The IP address of the SorryServer. This will be used when all servers under this VIP are down.
* `vip` - The virtual IP address.

---

A `server` block exports the following:

* `check_path` - The path used when checking by HTTP/HTTPS.
* `check_protocol` - The protocol used for health checks. This will be one of [`http`/`https`/`tcp`/`ping`].
* `check_status` - The response-code to expect when checking by HTTP/HTTPS.
* `enabled` - The flag to enable as destination of load balancing.
* `ip_address` - The IP address of the destination server.



