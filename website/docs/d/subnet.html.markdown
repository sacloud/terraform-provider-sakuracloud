---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_subnet"
sidebar_current: "docs-sakuracloud-datasource-subnet"
description: |-
  Get information on a SakuraCloud Subnet.
---

# sakuracloud\_subnet

Use this data source to retrieve information about a SakuraCloud Subnet.

## Example Usage

```hcl
data "sakuracloud_subnet" "foobar" {
  internet_id = sakuracloud_internet.foobar.id
  index       = 0
}
```

## Argument Reference

 * `internet_id` - (Required) The ID of the Internet resource.
 * `index` - (Required) The index of the target subnet.

## Attributes Reference

* `id` - The ID of the resource.
* `nw_mask_len` - Network mask length of the Subnet.
* `next_hop` - Next hop address.
* `switch_id` - The ID of the switch connected to the Subnet.
* `nw_address` -  The network address.
* `min_ipaddress` - Min global IP address.
* `max_ipaddress` - Max global IP address.
* `ipaddresses` - Global IP address list.
* `zone` - The ID of the zone to which the resource belongs.
