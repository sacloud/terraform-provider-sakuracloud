---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_firewall"
sidebar_current: "docs-sakuracloud-resource-vpc-firewall"
description: |-
  Provides a SakuraCloud VPC Router Firewall resource. This can be used to create and delete VPC Router Firewalls.
---

# sakuracloud\_vpc\_router\_firewall

Provides a SakuraCloud VPC Router Firewall resource. This can be used to create and delete VPC Router Firewalls.

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

# Create a new Firewall config
resource "sakuracloud_vpc_router_firewall" "send_fw" {
  vpc_router_id              = sakuracloud_vpc_router.foobar.id
  vpc_router_interface_index = 1

  direction = "send"
  expressions {
    protocol    = "tcp"
    source_nw   = ""
    source_port = "80"
    dest_nw     = ""
    dest_port   = ""
    allow       = true
    logging     = true
    description = "allow http"
  }

  expressions {
    protocol    = "ip"
    source_nw   = ""
    source_port = ""
    dest_nw     = ""
    dest_port   = ""
    allow       = false
    logging     = false
    description = "deny all"
  }
}


```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `vpc_router_interface_index` - (Required) The NIC index of VPC Router.
* `direction` - (Required) Direction of filtering packets.  
Valid value is one of the following: [ "send" / "receive" ]
* `expressions` - (Required) Filtering rules. It contains some attributes to [Expressions](#expressions).
* `zone` - (Optional) The ID of the zone to which the resource belongs.

### Expressions

Attributes for Expressions:

* `protocol` - (Required) Protocol used in health check.  
Valid value is one of the following: [ "tcp" / "udp" / "icmp" / "ip" ]
* `source_nw` - (Required) Target source network IP address or CIDR or range.  
Valid format is one of the following:   
  * IP address: `"xxx.xxx.xxx.xxx"`
  * CIDR: `"xxx.xxx.xxx.xxx/nn"`
  * Range: `"xxx.xxx.xxx.xxx/yyy.yyy.yyy.yyy"`
* `source_port` - (Required) Target source port.
Valid format is one of the following:
  * Number: `"0"` - `"65535"`
  * Range: `"xx-yy"`
  * Range (hex): `"0xPPPP/0xMMMM"`
* `dest_nw` - (Required) Target destination network IP address or CIDR or range.  
  Valid format is one of the following:   
    * IP address: `"xxx.xxx.xxx.xxx"`
    * CIDR: `"xxx.xxx.xxx.xxx/nn"`
    * Range: `"xxx.xxx.xxx.xxx/yyy.yyy.yyy.yyy"`
* `dest_port` - (Required) Target destination port.
Valid format is one of the following:
  * Number: `"0"` - `"65535"`
  * Range: `"xx-yy"`
  * Range (hex): `"0xPPPP/0xMMMM"`
* `allow` - (Required) The flag for allow/deny packets.
* `logging` - (Required) The flag for enable/disable logging.
* `description` - (Optional) The description of the expression.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `direction` - Direction of filtering packets.
* `expressions` - Filtering rules. It contains some attributes to [Expressions](#expressions).
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of VPC Router Firewall is not supported.
