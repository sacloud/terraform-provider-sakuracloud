---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_dhcp_static_mapping"
sidebar_current: "docs-sakuracloud-resource-vpc-dhcpmapping"
description: |-
  Provides a SakuraCloud VPC Router DHCP Static Mapping resource. This can be used to create and delete VPC Router DHCP Static Mappings.
---

# sakuracloud\_vpc\_router\_dhcp\_static\_mapping

Provides a SakuraCloud VPC Router DHCP Static Mapping resource. This can be used to create and delete VPC Router DHCP Static Mappings.

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

# Create a new DHCP Static Mapping config
resource "sakuracloud_vpc_router_dhcp_static_mapping" "dhcp_map" {
  vpc_router_id             = sakuracloud_vpc_router.foobar.id
  vpc_router_dhcp_server_id = sakuracloud_vpc_router_dhcp_server.foobar.id

  macaddress = "aa:bb:cc:aa:bb:cc"
  ipaddress  = "192.168.2.51"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `vpc_router_dhcp_server_id` - (Required) The ID of VPC Router DHCP Server.
* `macaddress` - (Required) The IP address mapped by MAC address.
* `ipaddress` - (Required) The MAC address to be the key of the mapping. 
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `macaddress` - The IP address mapped by MAC address.
* `ipaddress` - The MAC address to be the key of the mapping. 
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router DHCP Static Mapping is not supported.
