---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_subnet"
sidebar_current: "docs-sakuracloud-resource-networking-subnet"
description: |-
  Provides a SakuraCloud Subnet resource. This can be used to create, update, and delete Subnets.
---

# sakuracloud\_subnet

Provides a SakuraCloud Subnet resource. This can be used to create, update, and delete Subnets.

## Example Usage

```hcl
# Create a new Subnet
resource "sakuracloud_subnet" "foobar" {
  name        = "foobar"
  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `internet_id` - (Required) The ID of the Internet resource.
* `nw_mask_len` - (Optional) Network mask length.  
Valid value is one of the following: [ 28 (default) / 27 / 26 ]
* `next_hop` - (Optional) The next hop IP address.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `nw_mask_len` - Network mask length of the Subnet.
* `next_hop` - The next hop IP address.
* `subnet_id` - The ID of the subnet connected to the Subnet.
* `nw_address` -  The network address.
* `min_ipaddress` - Min global IP address.
* `max_ipaddress` - Max global IP address.
* `ipaddresses` - Global IP address list.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Subnets can be imported using the Subnet ID.

```
$ terraform import sakuracloud_subnet.foobar <subnet_id>
```
