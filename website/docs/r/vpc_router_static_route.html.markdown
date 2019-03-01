---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_static_route"
sidebar_current: "docs-sakuracloud-resource-vpc-staticroute"
description: |-
  Provides a SakuraCloud VPC Router Static Route resource. This can be used to create and delete VPC Router Static Route.
---

# sakuracloud\_vpc\_router\_static\_route

Provides a SakuraCloud VPC Router Static Route resource. This can be used to create and delete VPC Router Static Route.

## Example Usage

```hcl
# Create a new VPC Router(standard)
resource "sakuracloud_vpc_router" "foobar" {
  name = "foobar"
}

# Add NIC to the VPC Router
resource "sakuracloud_vpc_router_interface" "eth1" {
  vpc_router_id = sakuracloud_vpc_router.foobar.id
  index         = 1
  switch_id     = sakuracloud_switch.foobar.id
  ipaddress     = ["192.168.2.1"]
  nw_mask_len   = 24
}

# Add Static Route config
resource "sakuracloud_vpc_router_static_route" "route" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id

  prefix   = "10.0.0.0/8"
  next_hop = "192.2.0.11"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `vpc_router_interface_id` - (Required) The ID of VPC Router Interface.
* `prefix` - (Required) The prefix of the Static Route.
* `next_hop` - (Required) The next hop IP address of the Static Route.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `prefix` - The prefix of the Static Route.
* `next_hop` - The next hop IP address of the Static Route.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router Static Route is not supported.
