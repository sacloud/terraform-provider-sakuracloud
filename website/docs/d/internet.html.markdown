---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
sidebar_current: "docs-sakuracloud-datasource-internet"
description: |-
  Get information on a SakuraCloud internet.
---

# sakuracloud\_internet

Use this data source to retrieve information about a SakuraCloud internet(switch+router).

## Example Usage

```hcl
data sakuracloud_internet "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The filter value list of name.
 * `tag_selectors` - (Optional) The filter value list of tags.
 * `filter` - (Optional) The map of filter key and value.
 * `zone` - (Optional) The ID of the zone.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `nw_mask_len` - Network mask length.
* `band_width` - Bandwidth of outbound traffic.
* `switch_id` - The ID of the switch.
* `server_ids` - The IDs of the server connected to the switch.
* `nw_address` - The network address.
* `gateway` - The network gateway address of the switch.
* `min_ipaddress` - Minimum global ip address.
* `max_ipaddress` - Maximum global ip address.
* `ipaddresses` - Global ip address list.
* `enable_ipv6` - The ipv6 enabled flag.
* `ipv6_prefix` - Address prefix of ipv6 network.
* `ipv6_prefix_len` - Address prefix length of ipv6 network.
* `ipv6_nw_address` - The ipv6 network address.
* `description` - The description of the resource.
* `tags` - The tag list of the resource.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
