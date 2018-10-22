---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
sidebar_current: "docs-sakuracloud-datasource-load-balancer"
description: |-
  Get information on a SakuraCloud Load Balancer.
---

# sakuracloud\_load\_balancer

Use this data source to retrieve information about a SakuraCloud Load Balancer.

## Example Usage

```hcl
data "sakuracloud_load_balancer" "foobar" {
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
* `switch_id` - The ID of the Switch connected to the Load Balancer.
* `vrid` - VRID used when high-availability mode enabled.
* `high_availability` - The flag of enable/disable high-availability mode.
* `plan` - The name of the resource plan. 
* `ipaddress1` - The primary IP address of the Load Balancer.
* `ipaddress2` - The secondly IP address of the Load Balancer. Used when high-availability mode enabled.
* `nw_mask_len` - Network mask length.
* `default_route` - Default gateway address of the Load Balancer.	 
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
