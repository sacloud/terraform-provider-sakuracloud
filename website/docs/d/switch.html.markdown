---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
sidebar_current: "docs-sakuracloud-datasource-switch"
description: |-
  Get information on a SakuraCloud switch.
---

# sakuracloud\_switch

Use this data source to retrieve information about a SakuraCloud switch.

## Example Usage

```hcl
data sakuracloud_switch "foobar" {
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
* `bridge_id` - The ID of the bridge connected to the switch.
* `name` - The name of the resource.
* `server_ids` - The IDs of the server connected to the switch.
* `description` - The description of the resource.
* `tags` - The tag list of the resource.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
