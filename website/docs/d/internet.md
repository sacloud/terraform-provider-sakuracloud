---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
subcategory: "Networking"
description: |-
  Get information about an existing Switch+Router.
---

# Data Source: sakuracloud_internet

Get information about an existing Switch+Router.

## Example Usage

```hcl
data "sakuracloud_internet" "foobar" {
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

* `id` - The id of the Switch+Router.
* `band_width` - The bandwidth of the network connected to the Internet in Mbps.
* `description` - The description of the Switch+Router.
* `enable_ipv6` - The flag to enable IPv6.
* `gateway` - The IP address of the gateway used by Switch+Router.
* `icon_id` - The icon id attached to the Switch+Router.
* `ip_addresses` - A list of assigned global address to the Switch+Router.
* `ipv6_network_address` - The IPv6 network address assigned to the Switch+Router.
* `ipv6_prefix` - The network prefix of assigned IPv6 addresses to the Switch+Router.
* `ipv6_prefix_len` - The bit length of IPv6 network prefix.
* `max_ip_address` - Maximum IP address in assigned global addresses to the Switch+Router.
* `min_ip_address` - Minimum IP address in assigned global addresses to the Switch+Router.
* `name` - The name of the Switch+Router.
* `netmask` - The bit length of the subnet assigned to the Switch+Router.
* `network_address` - The IPv4 network address assigned to the Switch+Router.
* `server_ids` - A list of the ID of Servers connected to the Switch+Router.
* `switch_id` - The id of the switch connected from the Switch+Router.
* `tags` - Any tags assigned to the Switch+Router.
* `zone` - The name of zone that the Switch+Router is in (e.g. `is1a`, `tk1a`).



