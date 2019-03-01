---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_internet"
sidebar_current: "docs-sakuracloud-datasource-internet"
description: |-
  Get information on a SakuraCloud Internet (Switch+Router).
---

# sakuracloud\_internet

Use this data source to retrieve information about a SakuraCloud Internet (Switch+Router).

## Example Usage

```hcl
data "sakuracloud_internet" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.
 * `zone` - (Optional) The ID of the zone.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `nw_mask_len` - Network mask length.
* `band_width` - Bandwidth of outbound traffic.
* `switch_id` - The ID of the switch.
* `server_ids` - The ID list of the servers connected to the switch.
* `nw_address` - The network address.
* `gateway` - The network gateway address of the switch.
* `min_ipaddress` - Min global IP address.
* `max_ipaddress` - Max global IP address.
* `ipaddresses` - Global IP address list.
* `enable_ipv6` - The ipv6 enabled flag.
* `ipv6_prefix` - Address prefix of ipv6 network.
* `ipv6_prefix_len` - Address prefix length of ipv6 network.
* `ipv6_nw_address` - The ipv6 network address.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
