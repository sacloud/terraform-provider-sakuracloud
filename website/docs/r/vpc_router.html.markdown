---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router"
sidebar_current: "docs-sakuracloud-resource-vpc-router"
description: |-
  Provides a SakuraCloud VPC Router resource. This can be used to create, update, and delete VPC Routers.
---

# sakuracloud\_vpc\_router

Provides a SakuraCloud VPC Router resource. This can be used to create, update, and delete VPC Routers.

## Example Usage

```hcl
# Create a new VPC Router(standard)
resource "sakuracloud_vpc_router" "foobar" {
  name = "foobar"

  #syslog_host         = "192.168.11.1"
  #internet_connection = true

  description = "description"
  tags        = ["foo", "bar"]
}

# Create a new VPC Router(premium or highspec)
resource "sakuracloud_vpc_router" "foobar1" {
  name       = "foobar"
  plan       = "premium"
  switch_id  = sakuracloud_internet.foobar.switch_id
  vip        = sakuracloud_internet.foobar.ipaddresses[0]
  ipaddress1 = sakuracloud_internet.foobar.ipaddresses[1]
  ipaddress2 = sakuracloud_internet.foobar.ipaddresses[2]
  #aliases   = [sakuracloud_internet.foobar.ipaddresses[3]] 
  vrid = 1
  
  interface {
    switch_id   = sakuracloud_switch.sw.id
    vip         = "192.168.11.1"
    ipaddress   = ["192.168.11.2" , "192.168.11.3"]
    nw_mask_len = 24 
  }

  port_forwarding {
    protocol        = "udp"
    global_port     = 10022
    private_address = "192.168.11.11"
    private_port    = 22
    description     = "desc"
  }

  static_nat {
    global_address  = sakuracloud_internet.router1.ipaddresses[3]
    private_address = "192.168.11.12"
    description     = "desc"
  }

  firewall {
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
        description = "desc"
    }

    expressions {
        protocol    = "ip"
        source_nw   = ""
        source_port = ""
        dest_nw     = ""
        dest_port   = ""
        allow       = false
        logging     = true
        description = "desc"
    }
  }
  
  dhcp_server {
    vpc_router_interface_index = 1

    range_start = "192.168.11.11"
    range_stop  = "192.168.11.20"
    dns_servers = ["8.8.8.8", "8.8.4.4"]
  }

  dhcp_static_mapping {
    ipaddress  = "192.168.11.10"
    macaddress = "aa:bb:cc:aa:bb:cc"
  }

  l2tp {
    pre_shared_secret = "example"
    range_start       = "192.168.11.21"
    range_stop        = "192.168.11.30"
  }

  pptp {
    range_start = "192.168.11.31"
    range_stop  = "192.168.11.40"
  }

  user {
    name     = "username"
    password = "password"
  }
 
  site_to_site_vpn {
    peer              = "8.8.8.8"
    remote_id         = "8.8.8.8"
    pre_shared_secret = "example"
    routes            = ["10.0.0.0/8"]
    local_prefix      = ["192.168.21.0/24"]
  }

  static_route {
    prefix   = "172.16.0.0/16"
    next_hop = "192.168.11.99"
  }  
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `plan` - (Optional) The plan of the VPC Router.   
Valid value is one of the following: [ "standard" (default) / "premium" / "highspec" / "highspec4000" ]
* `switch_id` - (Required) The ID of the switch connected to the VPC Router. Used when plan is `premium` or `highspec` or `highspec4000`.
* `vrid` - (Required) VRID used when plan is `premium` or `highspec` or `highspec4000`.
* `ipaddress1` - (Required) The primary IP address of the VPC Router.
* `ipaddress2` - (Optional) The secondly IP address of the VPC Router. Used when plan is `premium` or `highspec` or `highspec4000`.
* `vip` - (Optional) The Virtual IP address of the VPC Router. Used when plan is `premium` or `highspec` or `highspec4000`.
* `aliases` - (Optional) The IP address aliase list. Used when plan is `premium` or `highspec` or `highspec4000`.
* `interface` - (Optional) The private NICs. It contains some attributes to [interface](#interface). 
* `port_forwarding` - (Optional) The port forwarding settings. It contains some attributes to [port_forwarding](#port_forwarding). 
* `static_nat` - (Optional) The static NAT settings. It contains some attributes to [static_nat](#static_nat).  
* `firewall` - (Optional) The firewall settings. It contains some attributes to [firewall](#firewall).  
* `dhcp_server` - (Optional) The DHCP server settings. It contains some attributes to [dhcp_server](#dhcp_server).   
* `dhcp_static_mapping` - (Optional) The DHCP static mapping settings. It contains some attributes to [dhcp_static_mapping](#dhcp_static_mapping).  
* `l2pt` - (Optional) The L2TP/IPsec settings. It contains some attributes to [l2tp](#l2tp).  
* `pptp` - (Optional) The PPTP settings. It contains some attributes to [pptp](#pptp).  
* `user` - (Optional) The remote access user settings. It contains some attributes to [user](#user).  
* `site_to_site_vpn` - (Optional) The Site-to-Site VPN settings. It contains some attributes to [site_to_site_vpn](#site_to_site_vpn).  
* `static_route` - (Optional) The static route settings. It contains some attributes to [static_route](#static_route).  
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the VPC Router.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `plan` - The name of the resource plan. 
* `switch_id` - The ID of the Switch connected to the VPC Router (eth0).
* `vip` - Virtual IP address of the VPC Router. Used when plan is in `premium` or `highspec` or `highspec4000`.
* `ipaddress1` - The primary IP address of the VPC Router.
* `ipaddress2` - The secondly IP address of the VPC Router. Used when plan is in `premium` or `highspec` or `highspec4000`.
* `vrid` - VRID used when plan is in `premium` or `highspec` or `highspec4000`.
* `aliases` - The IP address aliase list. Used when plan is in `premium` or `highspec` or `highspec4000`.
* `global_address` - Global IP address of the VPC Router.
* `syslog_host` - The destination HostName/IP address to send log.	
* `internet_connection` - The flag of enable/disable connection from the VPC Router to the Internet.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

### `interface`

The following arguments are supported:

* `index` - (Required) The NIC index of VPC Router Interface.
* `switch_id` - (Required) The ID of the switch connected to the VPC Router.
* `ipaddresses` - (Required) The IP address list of the VPC Router Interfaces.
* `vip` - (Optional) The Virtual IP address of the VPC Router Interface. Used when VPC Router's plan is `premium` or `highspec` or `highspec4000`.
* `nw_mask_len` - (Optional) Network mask length of the VPC Router Interface.

### `port_forwarding`

The following arguments are supported:

* `protocol` - (Required) The target protocol of the Port Forwarding.  
Valid value is one of the following: [ "tcp" (default) / "udp" ]
* `global_port` - (Required) The global port of the Port Forwarding.
* `private_address` - (Required) The destination private IP address of the Port Forwarding.
* `private_port` - (Required) The destination port number of the Port Forwarding.
* `description` - (Optional) The description of the resource.

### `static_nat`

The following arguments are supported:

* `global_address` - (Required) The global IP address of the Static NAT.
* `private_address` - (Required) The private IP address of the Static NAT.
* `description` - (Optional) The description of the resource.

### `firewall`

The following arguments are supported:

* `vpc_router_interface_index` - (Required) The NIC index of VPC Router.
* `direction` - (Required) Direction of filtering packets.  
Valid value is one of the following: [ "send" / "receive" ]
* `expressions` - (Required) Filtering rules. It contains some attributes to [Expressions](#expressions).

#### Expressions

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

### `dhcp_server`

The following arguments are supported:

* `vpc_router_interface_index` - (Required) The NIC index of VPC Router running DHCP Server.
* `range_start` - (Required) Start IP address of address range to be assigned by DHCP.
* `range_stop` - (Required) End IP address of address range to be assigned by DHCP.
* `dns_servers` - (Required) DNS server list to be assigned by DHCP.  

### `dhcp_static_mapping`

The following arguments are supported:

* `macaddress` - (Required) The IP address mapped by MAC address.
* `ipaddress` - (Required) The MAC address to be the key of the mapping. 

### `l2pt`

The following arguments are supported:

* `pre_shared_secret` - (Required) The pre shared secret for L2TP.
* `range_start` - (Required) Start IP address of address range to be assigned by L2TP.
* `range_stop` - (Required) End IP address of address range to be assigned by L2TP.

### `pptp`

The following arguments are supported:

* `range_start` - (Required) Start IP address of address range to be assigned by PPTP.
* `range_stop` - (Required) End IP address of address range to be assigned by PPTP.

### `user`

The following arguments are supported:

* `name` - (Required) The user name.
* `password` - (Required) The password.

### `site_to_site_vpn`

The following arguments are supported:

* `peer` - (Required) The peer IP address.
* `remote_id` - (Required) The IPSec ID of target.
* `pre_shared_secret` - (Required) The pre shared secret for IPSec.
* `routes` - (Required) The routing prefix.
* `local_prefix` - (Required) The local prefix.

#### Attributes Reference

The following attributes are exported:

* `peer` - The peer IP address.
* `remote_id` - The IPSec ID of target.
* `pre_shared_secret` - The pre shared secret for IPSec.
* `routes` - The routing prefix.
* `local_prefix` - The local prefix.
* `esp_authentication_protocol` - ESP authentication protocol.
* `esp_dh_group` - ESP DH group.
* `esp_encryption_protocol` - ESP encryption protocol.
* `esp_lifetime` - ESP lifetime.
* `esp_mode` - ESP mode.
* `esp_perfect_forward_secrecy` - ESP perfect forward secrecy.
* `ike_authentication_protocol` - IKE authentication protocol.
* `ike_encryption_protocol` - IKE encryption protocol.
* `ike_lifetime` - IKE lifetime.
* `ike_mode` - IKE mode.
* `ike_perfect_forward_secrecy` - IKE perfect forward secrecy.
* `ike_pre_shared_secret` - IKE pre shared secret.
* `peer_id` - Peer ID.
* `peer_inside_networks` - Peer inside networks.
* `peer_outside_ipaddress` - Peer outsite ipaddress.
* `vpc_router_inside_networks` - VPC Router inside networks.
* `vpc_router_outside_ipaddress` - VPC Router outside IP address.
* `zone` - The ID of the zone to which the resource belongs.

### `static_route`

The following arguments are supported:

* `prefix` - (Required) The prefix of the Static Route.
* `next_hop` - (Required) The next hop IP address of the Static Route.

## Import

VPC Routers can be imported using the VPC Router ID.

```
$ terraform import sakuracloud_vpc_router.foobar <vpc_router_id>
```
