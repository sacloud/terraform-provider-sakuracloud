---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_switch"
sidebar_current: "docs-sakuracloud-datasource-switch"
description: |-
  Get information on a SakuraCloud Switch.
---

# sakuracloud\_switch

Use this data source to retrieve information about a SakuraCloud Switch.

## Example Usage

```hcl
data "sakuracloud_switch" "foobar" {
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
* `bridge_id` - The ID of the bridge connected to the switch.
* `name` - The name of the resource.
* `server_ids` - The ID list of the servers connected to the switch.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
