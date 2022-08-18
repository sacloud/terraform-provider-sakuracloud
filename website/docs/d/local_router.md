---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_local_router"
subcategory: "Networking"
description: |-
  Get information about an existing Local Router.
---

# Data Source: sakuracloud_local_router

Get information about an existing Local Router.

## Example Usage

```hcl
data "sakuracloud_local_router" "foobar" {
  filter {
    names = ["foobar"]
  }
}
```

## Argument Reference

* `filter` - (Optional) One or more values used for filtering, as defined below.

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

* `id` - The id of the Local Router.
* `description` - The description of the LocalRouter.
* `icon_id` - The icon id attached to the LocalRouter.
* `name` - The name of the LocalRouter.
* `network_interface` - A list of `network_interface` blocks as defined below.
* `peer` - A list of `peer` blocks as defined below.
* `secret_keys` - A list of secret key used for peering from other LocalRouters.
* `static_route` - A list of `static_route` blocks as defined below.
* `switch` - A list of `switch` blocks as defined below.
* `tags` - Any tags assigned to the LocalRouter.


---

A `network_interface` block exports the following:

* `ip_addresses` - The list of IP address assigned to the LocalRouter.
* `netmask` - The bit length of the subnet assigned to the LocalRouter.
* `vip` - The virtual IP address.
* `vrid` - The Virtual Router Identifier.

---

A `peer` block exports the following:

* `description` - The description of the LocalRouter.
* `enabled` - The flag to enable the LocalRouter.
* `peer_id` - The ID of the peer LocalRouter.
* `secret_key` - The secret key of the peer LocalRouter.

---

A `static_route` block exports the following:

* `next_hop` - The IP address of the next hop.
* `prefix` - The CIDR block of destination.

---

A `switch` block exports the following:

* `category` - The category name of connected services (e.g. `cloud`, `vps`).
* `code` - The resource ID of the Switch.
* `zone_id` - The id of the Zone.


