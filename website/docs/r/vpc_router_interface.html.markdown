---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_interface"
sidebar_current: "docs-sakuracloud-resource-vpc-interface"
description: |-
  Provides a SakuraCloud VPC Router Interface resource. This can be used to create and delete VPC Router Interfaces.
---

# sakuracloud\_vpc\_router\_interface

Provides a SakuraCloud VPC Router Interface resource. This can be used to create and delete VPC Router Interfaces.

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

```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `index` - (Required) The NIC index of VPC Router Interface.
* `switch_id` - (Required) The ID of the switch connected to the VPC Router.
* `ipaddresses` - (Required) The IP address list of the VPC Router Interfaces.
* `vip` - (Optional) The Virtual IP address of the VPC Router Interface. Used when VPC Router's plan is `premium` or `highspec`.
* `nw_mask_len` - (Optional) Network mask length of the VPC Router Interface.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the VPC Router.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `vpc_router_id` - The ID of the Internet resource.
* `index` - The NIC index of VPC Router Interface.
* `switch_id` - The ID of the switch connected to the VPC Router.
* `ipaddresses` - The IP address list of the VPC Router Interfaces.
* `vip` - The Virtual IP address of the VPC Router Interface. Used when VPC Router's plan is `premium` or `highspec`.
* `nw_mask_len` - Network mask length of the VPC Router Interface.
* `graceful_shutdown_timeout` - The wait time (seconds) to do graceful shutdown the VPC Router.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router Interface is not supported.
