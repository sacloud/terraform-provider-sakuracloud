---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_dhcp_server"
sidebar_current: "docs-sakuracloud-resource-vpc-dhcpserver"
description: |-
  Provides a SakuraCloud VPC Router DHCP Server resource. This can be used to create, update, and delete VPC Router DHCP Servers.
---

# sakuracloud\_vpc\_router\_dhcp\_server

Provides a SakuraCloud VPC Router DHCP Server resource. This can be used to create, update, and delete VPC Router DHCP Servers.

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

# Create a new DHCP Server under the VPC Router
resource "sakuracloud_vpc_router_dhcp_server" "foobar" {
  vpc_router_id              = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_index = sakuracloud_vpc_router_interface.eth1.index

  range_start = "192.168.2.101"
  range_stop  = "192.168.2.200"
  dns_servers = ["8.8.8.8", "8.8.4.4"]
}
```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `vpc_router_interface_index` - (Required) The NIC index of VPC Router running DHCP Server.
* `range_start` - (Required) Start IP address of address range to be assigned by DHCP.
* `range_stop` - (Required) End IP address of address range to be assigned by DHCP.
* `dns_servers` - (Required) DNS server list to be assigned by DHCP.  
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `range_start` - Start IP address of address range to be assigned by DHCP.
* `range_stop` - End IP address of address range to be assigned by DHCP.
* `dns_servers` - DNS server list to be assigned by DHCP.  
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router DHCP Server is not supported.
