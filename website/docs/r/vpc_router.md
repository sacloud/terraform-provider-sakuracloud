---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud VPC Router.
---

# sakuracloud_vpc_router

Manages a SakuraCloud VPC Router.

## Argument Reference

* `aliases` - (Optional) A list of ip alias to assign to the VPC Router. This can only be specified if `plan` is not `standard`.
* `description` - (Optional) The description of the VPCRouter. The length of this value must be in the range [`1`-`512`].
* `dhcp_server` - (Optional) One or more `dhcp_server` blocks as defined below.
* `dhcp_static_mapping` - (Optional) One or more `dhcp_static_mapping` blocks as defined below.
* `firewall` - (Optional) One or more `firewall` blocks as defined below.
* `icon_id` - (Optional) The icon id to attach to the VPCRouter.
* `internet_connection` - (Optional) The flag to enable connecting to the Internet from the VPC Router. Default:`true`.
* `ip_addresses` - (Optional) The list of the IP address to assign to the VPC Router. This is required only one value when `plan` is `standard`, two values otherwise. Changing this forces a new resource to be created.
* `l2tp` - (Optional) A `l2tp` block as defined below.
* `name` - (Required) The name of the VPCRouter. The length of this value must be in the range [`1`-`64`].
* `network_interface` - (Optional) A list of additional network interface setting. This doesn't include primary network interface setting.
* `plan` - (Optional) The plan name of the VPCRouter. This must be one of [`standard`/`premium`/`highspec`/`highspec4000`]. Changing this forces a new resource to be created. Default:`standard`.
* `port_forwarding` - (Optional) One or more `port_forwarding` blocks as defined below.
* `pptp` - (Optional) A `pptp` block as defined below.
* `site_to_site_vpn` - (Optional) One or more `site_to_site_vpn` blocks as defined below.
* `static_nat` - (Optional) One or more `static_nat` blocks as defined below.
* `static_route` - (Optional) One or more `static_route` blocks as defined below.
* `switch_id` - (Optional) The id of the switch to connect. This is only required when when `plan` is not `standard`%!(EXTRA string=VPCRouter). Changing this forces a new resource to be created.
* `syslog_host` - (Optional) The ip address of the syslog host to which the VPC Router sends logs.
* `tags` - (Optional) Any tags to assign to the VPCRouter.
* `user` - (Optional) One or more `user` blocks as defined below.
* `vip` - (Optional) The virtual IP address of the VPC Router. This is only required when `plan` is not `standard`. Changing this forces a new resource to be created.
* `vrid` - (Optional) The Virtual Router Identifier. This is only required when `plan` is not `standard`. Changing this forces a new resource to be created.
* `zone` - (Optional) The name of zone that the VPCRouter will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


---

A `dhcp_server` block supports the following:

* `dns_servers` - (Optional) A list of IP address of DNS server to assign to DHCP client.
* `interface_index` - (Required) The index of the network interface on which to enable the DHCP service. This must be in the range [`1`-`7`].
* `range_start` - (Required) The start value of IP address range to assign to DHCP client.
* `range_stop` - (Required) The end value of IP address range to assign to DHCP client.

---

A `dhcp_static_mapping` block supports the following:

* `ip_address` - (Required) The static IP address to assign to DHCP client.
* `mac_address` - (Required) The source MAC address of static mapping.

---

A `firewall` block supports the following:

* `direction` - (Required) The direction to apply the firewall. This must be one of [`send`/`receive`].
* `expression` - (Required) One or more `expression` blocks as defined below.
* `interface_index` - (Optional) The index of the network interface on which to enable filtering. This must be in the range [`0`-`7`].

---

A `expression` block supports the following:

* `allow` - (Required) The flag to allow the packet through the filter.
* `description` - (Optional) The description of the expression. The length of this value must be in the range [`0`-`512`].
* `destination_network` - (Optional) A destination IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `destination_port` - (Optional) A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`.
* `logging` - (Optional) The flag to enable packet logging when matching the expression.
* `protocol` - (Required) The protocol used for filtering. This must be one of [`tcp`/`udp`/`icmp`/`ip`].
* `source_network` - (Optional) A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `source_port` - (Optional) A source port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`.

---

A `l2tp` block supports the following:

* `pre_shared_secret` - (Required) The pre shared secret for L2TP/IPsec.
* `range_start` - (Required) The start value of IP address range to assign to L2TP/IPsec client.
* `range_stop` - (Required) The end value of IP address range to assign to L2TP/IPsec client.

---

A `network_interface` block supports the following:

* `index` - (Required) The index of the network interface. This must be in the range [`1`-`7`].
* `ip_addresses` - (Required) A list of ip address to assign to the network interface. This is required only one value when `plan` is `standard`, two values otherwise.
* `netmask` - (Required) The bit length of the subnet to assign to the network interface.
* `switch_id` - (Required) The id of the connected switch.
* `vip` - (Optional) The virtual IP address to assign to the network interface. This is only required when `plan` is not `standard`.

---

A `port_forwarding` block supports the following:

* `description` - (Optional) The description of the port forwarding. The length of this value must be in the range [`0`-`512`].
* `private_ip` - (Required) The destination ip address of the port forwarding.
* `private_port` - (Required) The destination port number of the port forwarding. This will be a port number on a private network.
* `protocol` - (Required) The protocol used for port forwarding. This must be one of [`tcp`/`udp`].
* `public_port` - (Required) The source port number of the port forwarding. This must be a port number on a public network.

---

A `pptp` block supports the following:

* `range_start` - (Required) The start value of IP address range to assign to PPTP client.
* `range_stop` - (Required) The end value of IP address range to assign to PPTP client.

---

A `site_to_site_vpn` block supports the following:

* `local_prefix` - (Required) A list of CIDR block of the network under the VPC Router.
* `peer` - (Required) The IP address of the opposing appliance connected to the VPC Router.
* `pre_shared_secret` - (Required) The pre shared secret for the VPN. The length of this value must be in the range [`0`-`40`].
* `remote_id` - (Required) The id of the opposing appliance connected to the VPC Router. This is typically set same as value of `peer`.
* `routes` - (Required) A list of CIDR block of VPN connected networks.

---

A `static_nat` block supports the following:

* `description` - (Optional) The description of the static nat. The length of this value must be in the range [`0`-`512`].
* `private_ip` - (Required) The private IP address used for the static NAT.
* `public_ip` - (Required) The public IP address used for the static NAT.

---

A `static_route` block supports the following:

* `next_hop` - (Required) The IP address of the next hop.
* `prefix` - (Required) The CIDR block of destination.

---

A `user` block supports the following:

* `name` - (Required) The user name used to authenticate remote access.
* `password` - (Required) The password used to authenticate remote access.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the VPC Router


* `update` - (Defaults to 60 minutes) Used when updating the VPC Router

* `delete` - (Defaults to 20 minutes) Used when deregistering VPC Router



## Attribute Reference

* `id` - The id of the VPC Router.
* `public_ip` - The public ip address of the VPC Router.




