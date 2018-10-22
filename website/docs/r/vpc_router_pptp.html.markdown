---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_pptp"
sidebar_current: "docs-sakuracloud-resource-vpc-pptp"
description: |-
  Provides a SakuraCloud VPC Router PPTP resource. This can be used to create, update, and delete VPC Router PPTP.
---

# sakuracloud\_vpc\_router\_pptp

Provides a SakuraCloud VPC Router PPTP resource. This can be used to create, update, and delete VPC Router PPTP.

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

# Create a new PPTP config.
resource "sakuracloud_vpc_router_pptp" "pptp" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id

  range_start = "192.168.2.51"
  range_stop  = "192.168.2.100"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `vpc_router_interface_id` - (Required) The ID of VPC Router Interface.
* `range_start` - (Required) Start IP address of address range to be assigned by PPTP.
* `range_stop` - (Required) End IP address of address range to be assigned by PPTP.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `pre_shared_secret` - The pre shared secret for PPTP.
* `range_start` - Start IP address of address range to be assigned by PPTP.
* `range_stop` - End IP address of address range to be assigned by PPTP.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router PPTP is not supported.
