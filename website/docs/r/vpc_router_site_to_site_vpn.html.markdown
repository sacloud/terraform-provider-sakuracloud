---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router_site_to_site_vpn"
sidebar_current: "docs-sakuracloud-resource-vpc-s2s"
description: |-
  Provides a SakuraCloud VPC Router Site To Site VPN resource. This can be used to create and delete VPC Router Site To Site VPN.
---

# sakuracloud\_vpc\_router\_site\_to\_site\_vpn

Provides a SakuraCloud VPC Router Site To Site VPN resource. This can be used to create and delete VPC Router Site To Site VPN.

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

# Add Site to Site VPN config
resource "sakuracloud_vpc_router_site_to_site_vpn" "s2s" {
  vpc_router_id     = sakuracloud_vpc_router.foobar.id
  peer              = "172.16.0.1"
  remote_id         = "172.16.0.1"
  pre_shared_secret = "<your-pre-shared-secret>"
  routes            = ["10.0.0.0/8"]
  local_prefix      = ["192.168.2.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `vpc_router_id` - (Required) The ID of the Internet resource.
* `peer` - (Required) The peer IP address.
* `remote_id` - (Required) The IPSec ID of target.
* `pre_shared_secret` - (Required) The pre shared secret for IPSec.
* `routes` - (Required) The routing prefix.
* `local_prefix` - (Required) The local prefix.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
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

## Import (not supported)

Import of VPC Router Site To Site VPN is not supported.
