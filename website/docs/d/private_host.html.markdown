---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_private_host"
sidebar_current: "docs-sakuracloud-datasource-private-host"
description: |-
  Get information on a SakuraCloud Private Host.
---

# sakuracloud\_private\_host

Use this data source to retrieve information about a SakuraCloud Private Host.

## Example Usage

```hcl
data "sakuracloud_private_host" "foobar" {
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
* `hostname` - The HostName of the resource.
* `assigned_core` - The number of cores assigned to the Server.
* `assigned_memory` - The size of memory allocated to the Server (unit:`GB`).
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `zone` - The ID of the zone to which the resource belongs.
