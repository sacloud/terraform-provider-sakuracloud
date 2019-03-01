---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_nfs"
sidebar_current: "docs-sakuracloud-datasource-nfs"
description: |-
  Get information on a SakuraCloud NFS.
---

# sakuracloud\_nfs

Use this data source to retrieve information about a SakuraCloud NFS.

## Example Usage

```hcl
data "sakuracloud_nfs" "foobar" {
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
* `switch_id` - The ID of the Switch connected to the NFS.
* `plan` - The name of the resource plan.
* `ipaddress` - The IP address of the NFS.
* `nw_mask_len` - Network mask length.
* `default_route` - Default gateway address of the NFS.	 
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.


