---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
subcategory: "Appliance"
description: |-
  Get information about an existing Load Balancer.
---

# Data Source: sakuracloud_load_balancer

Get information about an existing Load Balancer.

## Example Usage

```hcl
data "sakuracloud_load_balancer" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.
* `zone` - (Optional) The name of zone that the LoadBalancer is in (e.g. `is1a`, `tk1a`).

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
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the Load Balancer.
* `description` - The description of the LoadBalancer.
* `network_interface` - A list of `network_interface` blocks as defined below.
* `icon_id` - The icon id attached to the LoadBalancer.
* `name` - The name of the LoadBalancer.
* `plan` - The plan name of the LoadBalancer. This will be one of [`standard`/`highspec`].
* `tags` - Any tags assigned to the LoadBalancer.
* `vip` - A list of `vip` blocks as defined below.

---

A `network_interface` block exports the following:

* `gateway` - The IP address of the gateway used by LoadBalancer.
* `ip_addresses` - The list of IP address assigned to the LoadBalancer.
* `netmask` - The bit length of the subnet assigned to the LoadBalancer.
* `switch_id` - The id of the switch connected from the LoadBalancer.
* `vrid` - The Virtual Router Identifier.

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

* `enabled` - The flag to enable as destination of load balancing.
* `ip_address` - The IP address of the destination server.
* `path` - The path used when checking by HTTP/HTTPS.
* `protocol` - The protocol used for health checks. This will be one of [`http`/`https`/`tcp`/`ping`].
* `status` - The response code to expect when checking by HTTP/HTTPS.


