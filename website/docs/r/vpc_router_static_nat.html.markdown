---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_static_nat"
sidebar_current: "docs-sakuracloud-resource-vpc-staticnat"
description: |-
  Provides a SakuraCloud VPC Router Static NAT resource. This can be used to create and delete VPC Router Static NAT.
---

# sakuracloud\_vpc\_router\_static\_nat

Provides a SakuraCloud VPC Router Static NAT resource. This can be used to create and delete VPC Router Static NAT.

## Example Usage

```hcl
# Create a new VPC Router(premium or highspec)
resource "sakuracloud_vpc_router" "foobar" {
  name       = "foobar"
  plan       = "premium"
  switch_id  = sakuracloud_internet.foobar.switch_id
  vip        = sakuracloud_internet.foobar.ipaddresses[0]
  ipaddress1 = sakuracloud_internet.foobar.ipaddresses[1]
  ipaddress2 = sakuracloud_internet.foobar.ipaddresses[2]
  #aliases   = [sakuracloud_internet.foobar.ipaddresses[3]] 
  vrid = 1
}

# Add NIC to the VPC Router
resource "sakuracloud_vpc_router_interface" "eth1" {
  vpc_router_id = sakuracloud_vpc_router.foobar.id
  index         = 1
  switch_id     = sakuracloud_switch.foobar.id
  vip           = "192.2.0.1"
  ipaddress     = ["192.2.0.2", "192.2.0.3"]
  nw_mask_len   = 24
}

# Add Static NAT config
resource "sakuracloud_vpc_router_static_nat" "snat" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id

  global_address  = sakuracloud_internet.router1.ipaddresses[3]
  private_address = "192.2.0.11"
  description     = "description"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `vpc_router_interface_id` - (Required) The ID of VPC Router Interface.
* `global_address` - (Required) The global IP address of the Static NAT.
* `private_address` - (Required) The private IP address of the Static NAT.
* `description` - (Optional) The description of the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `global_address` - The global IP address of the Static NAT.
* `private_address` - The private IP address of the Static NAT.
* `description` - The description of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router Static NAT is not supported.
