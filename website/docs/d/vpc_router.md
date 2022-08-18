---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router"
subcategory: "Appliance"
description: |-
  Get information about an existing VPC Router.
---

# Data Source: sakuracloud_vpc_router

Get information about an existing VPC Router.

## Example Usage

```hcl
data "sakuracloud_vpc_router" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```
## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.
* `zone` - (Optional) The name of zone that the VPC Router is in (e.g. `is1a`, `tk1a`).

---

A `filter` block supports the following:

* `condition` - (Optional) One or more name/values pairs used for filtering. There are several valid keys, for a full reference, check out finding section in the [SakuraCloud API reference](https://developer.sakura.ad.jp/cloud/api/1.1/).
* `id` - (Optional) The resource id on SakuraCloud used for filtering.
* `names` - (Optional) The resource names on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.
* `tags` - (Optional) The resource tags on SakuraCloud used for filtering. If multiple values ​​are specified, they combined as AND condition.

---

A `condition` block supports the following:

* `name` - (Required) The name of the target field. This value is case-sensitive.
* `values` - (Required) The values of the condition. If multiple values ​​are specified, they combined as AND condition.
* `operator` - (Optional) The filtering operator. This must be one of following: `partial_match_and`/`exact_match_or`. Default: `partial_match_and`


## Attribute Reference

* `id` - The id of the VPC Router.
* `description` - The description of the VPCRouter.
* `dhcp_server` - A list of `dhcp_server` blocks as defined below.
* `dhcp_static_mapping` - A list of `dhcp_static_mapping` blocks as defined below.
* `dns_forwarding` - A list of `dns_forwarding` blocks as defined below.
* `firewall` - A list of `firewall` blocks as defined below.
* `icon_id` - The icon id attached to the VPCRouter.
* `internet_connection` - The flag to enable connecting to the Internet from the VPC Router.
* `l2tp` - A list of `l2tp` blocks as defined below.
* `name` - The id of the switch connected from the VPCRouter.
* `plan` - The plan name of the VPCRouter. This will be one of [`standard`/`premium`/`highspec`/`highspec4000`].
* `port_forwarding` - A list of `port_forwarding` blocks as defined below. This represents a `Reverse NAT`.
* `pptp` - A list of `pptp` blocks as defined below.
* `private_network_interface` - A list of additional network interface setting. This doesn't include primary network interface setting.
* `public_ip` - The public ip address of the VPC Router.
* `public_netmask` - The bit length of the subnet to assign to the public network interface.
* `public_network_interface` - A list of additional network interface setting. This doesn't include primary network interface setting.
* `site_to_site_vpn` - A list of `site_to_site_vpn` blocks as defined below.
* `site_to_site_vpn_parameter` - A list of `site_to_site_vpn_parameter` blocks as defined below.
* `static_nat` - A list of `static_nat` blocks as defined below. This represents a `1:1 NAT`, doing static mapping to both send/receive to/from the Internet. This is only used when `plan` is not `standard`.
* `static_route` - A list of `static_route` blocks as defined below.
* `syslog_host` - The ip address of the syslog host to which the VPC Router sends logs.
* `tags` - Any tags assigned to the VPCRouter.
* `user` - A list of `user` blocks as defined below.
* `version` - The version of the VPC Router.
* `wire_guard` - A list of `wire_guard` blocks as defined below.


---

A `dhcp_server` block exports the following:

* `dns_servers` - A list of IP address of DNS server to assign to DHCP client.
* `interface_index` - The index of the network interface on which to enable the DHCP service. This will be between `1`-`7`.
* `range_start` - The start value of IP address range to assign to DHCP client.
* `range_stop` - The end value of IP address range to assign to DHCP client.

---

A `dhcp_static_mapping` block exports the following:

* `ip_address` - The static IP address to assign to DHCP client.
* `mac_address` - The source MAC address of static mapping.

---

A `dns_forwarding` block exports the following:

* `dns_servers` - A list of IP address of DNS server to forward to.
* `interface_index` - The index of the network interface on which to enable the DNS forwarding service.

---

A `firewall` block exports the following:

* `direction` - The direction to apply the firewall. This will be one of [`send`/`receive`].
* `expression` - A list of `expression` blocks as defined below.
* `interface_index` - The index of the network interface on which to enable filtering. This will be between `0`-`7`.

---

A `expression` block exports the following:

* `allow` - The flag to allow the packet through the filter.
* `description` - The description of the expression.
* `destination_network` - A destination IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `destination_port` - A destination port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`.
* `logging` - The flag to enable packet logging when matching the expression.
* `protocol` - The protocol used for filtering. This will be one of [`tcp`/`udp`/`icmp`/`ip`].
* `source_network` - A source IP address or CIDR block used for filtering (e.g. `192.0.2.1`, `192.0.2.0/24`).
* `source_port` - A source port number or port range used for filtering (e.g. `1024`, `1024-2048`). This is only used when `protocol` is `tcp` or `udp`.

---

A `l2tp` block exports the following:

* `pre_shared_secret` - The pre shared secret for L2TP/IPsec.
* `range_start` - The start value of IP address range to assign to L2TP/IPsec client.
* `range_stop` - The end value of IP address range to assign to L2TP/IPsec client.

---

A `port_forwarding` block exports the following:

* `description` - The description of the port forwarding.
* `private_ip` - The destination ip address of the port forwarding.
* `private_port` - The destination port number of the port forwarding. This will be a port number on a private network.
* `protocol` - The protocol used for port forwarding. This will be one of [`tcp`/`udp`].
* `public_port` - The source port number of the port forwarding. This will be a port number on a public network.

---

A `pptp` block exports the following:

* `range_start` - The start value of IP address range to assign to PPTP client.
* `range_stop` - The end value of IP address range to assign to PPTP client.

---

A `private_network_interface` block exports the following:

* `index` - The index of the network interface. This will be between `1`-`7`.
* `ip_addresses` - A list of ip address assigned to the network interface. This will be only one value when `plan` is `standard`, two values otherwise.
* `netmask` - The bit length of the subnet assigned to the network interface.
* `switch_id` - The id of the connected switch.
* `vip` - The virtual IP address assigned to the network interface. This is only used when `plan` is not `standard`.

---

A `public_network_interface` block exports the following:

* `aliases` - A list of ip alias assigned to the VPC Router. This is only used when `plan` is not `standard`.
* `ip_addresses` - The list of the IP address assigned to the VPC Router. This will be only one value when `plan` is `standard`, two values otherwise.
* `switch_id` - The id of the switch connected from the VPCRouter.
* `vip` - The virtual IP address of the VPC Router. This is only used when `plan` is not `standard`.
* `vrid` - The Virtual Router Identifier. This is only used when `plan` is not `standard`.

---

A `site_to_site_vpn` block exports the following:

* `local_prefix` - A list of CIDR block of the network under the VPC Router.
* `peer` - The IP address of the opposing appliance connected to the VPC Router.
* `pre_shared_secret` - The pre shared secret for the VPN.
* `remote_id` - The id of the opposing appliance connected to the VPC Router. This is typically set same as value of `peer`.
* `routes` - A list of CIDR block of VPN connected networks.

---

A `site_to_site_vpn_parameter` block exports the following:

* `encryption_algo` - .
* `esp` - A list of `esp` blocks as defined below.
* `hash_algo` - .
* `ike` - A list of `ike` blocks as defined below.

---

A `esp` block exports the following:

* `lifetime` - .

---

A `ike` block exports the following:

* `dpd` - A list of `dpd` blocks as defined below.
* `lifetime` - .

---

A `dpd` block exports the following:

* `interval` - .
* `timeout` - .

---

A `static_nat` block exports the following:

* `description` - The description of the static NAT.
* `private_ip` - The private IP address used for the static NAT.
* `public_ip` - The public IP address used for the static NAT.

---

A `static_route` block exports the following:

* `next_hop` - The IP address of the next hop.
* `prefix` - The CIDR block of destination.

---

A `user` block exports the following:

* `name` - The user name used to authenticate remote access.
* `password` - The password used to authenticate remote access.

---

A `wire_guard` block exports the following:

* `ip_address` - The IP address for WireGuard server.
* `peer` - A list of `peer` blocks as defined below.
* `public_key` - the public key of the WireGuard server.

---

A `peer` block exports the following:

* `ip_address` - The IP address for peer.
* `name` - the name of the peer.
* `public_key` - the public key of the WireGuard client.


