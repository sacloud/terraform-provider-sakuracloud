---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_zone"
subcategory: "Provider Data Sources"
description: |-
  Get information about an existing Zone.
---

# Data Source: sakuracloud_zone

Get information about an existing Zone.

## Example Usage

```hcl
data "sakuracloud_zone" "current" {}

data "sakuracloud_zone" "is1a" {
  name = "is1a"
}
```
## Argument Reference

* `name` - (Optional) The name of the zone (e.g. `is1a`,`tk1a`).



## Attribute Reference

* `id` - The id of the Zone.
* `description` - The description of the zone.
* `dns_servers` - A list of IP address of DNS server in the zone.
* `region_id` - The id of the region that the zone belongs.
* `region_name` - The name of the region that the zone belongs.
* `zone_id` - The id of the zone.



