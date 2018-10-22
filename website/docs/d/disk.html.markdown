---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_disk"
sidebar_current: "docs-sakuracloud-datasource-disk"
description: |-
  Get information on a SakuraCloud Disk.
---

# sakuracloud\_disk

Use this data source to retrieve information about a SakuraCloud Disk.

## Example Usage

```hcl
data "sakuracloud_disk" "foobar" {
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
* `plan` - The plan of the resource (`ssd`/`hdd`).
* `conector` - The disk connector of the resource (`virtio`/`ide`).
* `size` - Size of the resource (unit:`GB`).
* `server_id` - The ID of the server connected to the disk.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.
