---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_cdrom"
sidebar_current: "docs-sakuracloud-datasource-cdrom"
description: |-
  Get information on a SakuraCloud CD-ROM.
---

# sakuracloud\_cdrom

Use this data source to retrieve information about a SakuraCloud CD-ROM.

## Example Usage

```hcl
data "sakuracloud_cdrom" "foobar" {
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
* `size` - Size of the resource (unit:`GB`).
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
