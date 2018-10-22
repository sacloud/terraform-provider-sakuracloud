---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_port_forwarding"
sidebar_current: "docs-sakuracloud-resource-vpc-pf"
description: |-
  Provides a SakuraCloud VPC Router Port Forwarding resource. This can be used to create, update, and delete VPC Router Port Forwarding.
---

# sakuracloud\_vpc\_router\_port_forwarding

Provides a SakuraCloud VPC Router Port Forwarding resource. This can be used to create, update, and delete VPC Router Port Forwarding.

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

# Add Port Forwarding rule to the VPC Router
resource "sakuracloud_vpc_router_port_forwarding" "forward1" {
  vpc_router_id           = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_id = sakuracloud_vpc_router_interface.eth1.id

  protocol        = "tcp"
  global_port     = 10022
  private_address = "192.168.2.11"
  private_port    = 22
  description     = "description"
}


```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `vpc_router_interface_id` - (Required) The ID of VPC Router Interface.
* `protocol` - (Required) The target protocol of the Port Forwarding.  
Valid value is one of the following: [ "tcp" (default) / "udp" ]
* `global_port` - (Required) The global port of the Port Forwarding.
* `private_address` - (Required) The destination private IP address of the Port Forwarding.
* `private_port` - (Required) The destination port number of the Port Forwarding.
* `description` - (Optional) The description of the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `protocol` - The target protocol of the Port Forwarding.  
* `global_port` - The global port of the Port Forwarding.
* `private_address` - The destination private IP address of the Port Forwarding.
* `private_port` - The destination port number of the Port Forwarding.
* `description` - The description of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router Port Forwarding is not supported.
