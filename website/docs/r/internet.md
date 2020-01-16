---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Switch+Router.
---

# sakuracloud_internet

Manages a SakuraCloud Switch+Router.

## Example Usage

```hcl
resource "sakuracloud_internet" "foobar" {
  name = "foobar"

  netmask     = 28
  band_width  = 100
  enable_ipv6 = false

  description = "description"
  tags        = ["tag1", "tag2"]
}
```

## Argument Reference

* `name` - (Required) The name of the Switch+Router. The length of this value must be in the range [`1`-`64`].
* `band_width` - (Optional) The bandwidth of the network connected to the Internet in Mbps. `100`/`250`/`500`/`1000`/`1500`/`2000`/`2500`/`3000`/`5000`. Default:`100`.
* `netmask` - (Optional) The bit length of the subnet assigned to the Switch+Router. `26`/`27`/`28`. Changing this forces a new resource to be created. Default:`28`.
* `enable_ipv6` - (Optional) The flag to enable IPv6.

#### Common Arguments

* `description` - (Optional) The description of the Switch+Router. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the Switch+Router.
* `tags` - (Optional) Any tags to assign to the Switch+Router.
* `zone` - (Optional) The name of zone that the Switch+Router will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Switch+Router
* `update` - (Defaults to 60 minutes) Used when updating the Switch+Router
* `delete` - (Defaults to 20 minutes) Used when deleting Switch+Router

## Attribute Reference

* `id` - The id of the Switch+Router.
* `gateway` - The IP address of the gateway used by the Switch+Router.
* `ip_addresses` - A list of assigned global address to the Switch+Router.
* `ipv6_network_address` - The IPv6 network address assigned to the Switch+Router.
* `ipv6_prefix` - The network prefix of assigned IPv6 addresses to the Switch+Router.
* `ipv6_prefix_len` - The bit length of IPv6 network prefix.
* `max_ip_address` - Maximum IP address in assigned global addresses to the Switch+Router.
* `min_ip_address` - Minimum IP address in assigned global addresses to the Switch+Router.
* `network_address` - The IPv4 network address assigned to the Switch+Router.
* `server_ids` - A list of the ID of Servers connected to the Switch+Router.
* `switch_id` - The id of the switch.

