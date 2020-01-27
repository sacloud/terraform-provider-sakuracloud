---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_local_router"
subcategory: "Networking"
description: |-
  Manages a SakuraCloud Local Router.
---

# sakuracloud_local_router

Manages a SakuraCloud Local Router.

## Example Usage

```hcl
resource "sakuracloud_local_router" "foobar" {
  name        = "example"
  description = "descriptio"
  tags        = ["tag1", "tag2"]

  switch {
    code     = sakuracloud_switch.foobar.id
    category = "cloud"
    zone_id  = "is1a"
  }

  network_interface {
    vip          = "192.168.11.1"
    ip_addresses = ["192.168.11.11", "192.168.11.12"]
    netmask      = 24
    vrid         = 101
  }

  static_route {
    prefix   = "10.0.0.0/24"
    next_hop = "192.168.11.2"
  }
  static_route {
    prefix   = "172.16.0.0/16"
    next_hop = "192.168.11.3"
  }

  peer {
    peer_id     = data.sakuracloud_local_router.peer.id
    secret_key  = data.sakuracloud_local_router.peer.secret_keys[0]
    description = "description"
  }
}

resource "sakuracloud_switch" "foobar" {
  name = "example"
}

data "sakuracloud_local_router" "peer" {
  filter {
    names = ["peer"]
  }
}

```
## Argument Reference

* `name` - (Required) The name of the LocalRouter. The length of this value must be in the range [`1`-`64`].
* `network_interface` - (Required) An `network_interface` block as defined below.
* `switch` - (Required) A `switch` block as defined below.
* `peer` - (Optional) One or more `peer` blocks as defined below.
* `static_route` - (Optional) One or more `static_route` blocks as defined below.

---

A `network_interface` block supports the following:

* `ip_addresses` - (Required) A list of IP address to assign to the LocalRouter.
* `netmask` - (Required) The bit length of the subnet assigned to the LocalRouter. This must be in the range [`8`-`29`].
* `vip` - (Required) The virtual IP address.
* `vrid` - (Required) The Virtual Router Identifier.

---

A `peer` block supports the following:

* `peer_id` - (Required) The ID of the peer LocalRouter.
* `secret_key` - (Required) The secret key of the peer LocalRouter.
* `description` - (Optional) The description of the LocalRouter. The length of this value must be in the range [`1`-`512`].
* `enabled` - (Optional) The flag to enable the LocalRouter.

---

A `static_route` block supports the following:

* `next_hop` - (Required) The IP address of the next hop.
* `prefix` - (Required) The CIDR block of destination.

---

A `switch` block supports the following:

* `category` - (Optional) The category name of connected services (e.g. `cloud`, `vps`).
* `code` - (Required) The resource ID of the Switch.
* `zone_id` - (Required) The id of the Zone.

#### Common Arguments

* `description` - (Optional) The description of the LoadBalancer. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the LoadBalancer.
* `tags` - (Optional) Any tags to assign to the LoadBalancer.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 minutes) Used when creating the Local Router
* `update` - (Defaults to 20 minutes) Used when updating the Local Router
* `delete` - (Defaults to 20 minutes) Used when deleting Local Router

## Attribute Reference

* `id` - The id of the Local Router.
* `secret_keys` - A list of secret key used for peering from other LocalRouters.

