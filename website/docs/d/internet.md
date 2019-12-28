---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
subcategory: "Networking"
description: |-
  Get information about an existing Internet.
---

# Data Source: sakuracloud_internet

Get information about an existing Internet.

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

* `id` - The id of the Internet.
* `band_width` - The bandwidth of the network connected to the Internet in Mbps.
* `description` - The description of the switch+router.
* `enable_ipv6` - The flag to enable IPv6.
* `gateway` - The IP address of the gateway used by switch+router.
* `icon_id` - The icon id attached to the switch+router.
* `ip_addresses` - A list of assigned global address to the switch+router.
* `ipv6_network_address` - The IPv6 network address assigned to the switch+router.
* `ipv6_prefix` - The network prefix of assigned IPv6 addresses to the switch+router.
* `ipv6_prefix_len` - The bit length of IPv6 network prefix.
* `max_ip_address` - Maximum IP address in assigned global addresses to the switch+router.
* `min_ip_address` - Minimum IP address in assigned global addresses to the switch+router.
* `name` - The name of the switch+router.
* `netmask` - The bit length of the subnet assigned to the switch+router.
* `network_address` - The IPv4 network address assigned to the switch+router.
* `server_ids` - A list of the ID of Servers connected to the switch+router.
* `switch_id` - The id of the switch connected from the switch+router.
* `tags` - Any tags assigned to the switch+router.
* `zone` - The name of zone that the switch+router is in (e.g. `is1a`,`tk1a`).




