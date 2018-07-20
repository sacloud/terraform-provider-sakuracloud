---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_bridge"
sidebar_current: "docs-sakuracloud-datasource-bridge"
description: |-
  Get information on a SakuraCloud bridge.
---

# sakuracloud\_bridge

Use this data source to retrieve information about a SakuraCloud bridge.

## Example Usage

```hcl
data sakuracloud_bridge "foobar" {
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
* `switch_ids` - The IDs of the switch connected to the bridge.
* `description` - The description of the resource.
* `tags` - The tag list of the resource.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
