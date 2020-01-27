---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_zone"
subcategory: "Data Sources"
description: |-
  Get information on a SakuraCloud Zone.
---

# sakuracloud\_zone

Use this data source to retrieve information about a SakuraCloud Zone.

## Example Usage

```hcl
data sakuracloud_zone "current" {}

data sakuracloud_zone "is1a" {
  name = "is1a"
}
```

## Argument Reference

 * `name` - (Optional) The name of zone(default: use provider settings).

## Attributes Reference

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `zone_id` - The Id of the resource.
* `description` - The description of the resource.
* `region_id` - The ID of the region.
* `region_name` - The Name of the region.
* `dns_servers` - The IP Address list of the region.
