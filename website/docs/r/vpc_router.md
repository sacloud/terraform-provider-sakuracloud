---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud VPCRouter.
---

# sakuracloud_vpc_router

Manages a SakuraCloud VPCRouter.

## Argument Reference

* `aliases` - (Optional) .
* `description` - (Optional) .
* `dhcp_server` - (Optional) One or more `dhcp_server` blocks as defined below.
* `dhcp_static_mapping` - (Optional) One or more `dhcp_static_mapping` blocks as defined below.
* `firewall` - (Optional) One or more `firewall` blocks as defined below.
* `icon_id` - (Optional) .
* `internet_connection` - (Optional) . Defaults to `true`.
* `ip_addresses` - (Optional) . Changing this forces a new resource to be created.
* `l2tp` - (Optional) A `l2tp` block as defined below.
* `name` - (Required) .
* `network_interface` - (Optional) One or more `network_interface` blocks as defined below.
* `plan` - (Optional) . Changing this forces a new resource to be created. Defaults to `standard`.
* `port_forwarding` - (Optional) One or more `port_forwarding` blocks as defined below.
* `pptp` - (Optional) A `pptp` block as defined below.
* `site_to_site_vpn` - (Optional) One or more `site_to_site_vpn` blocks as defined below.
* `static_nat` - (Optional) One or more `static_nat` blocks as defined below.
* `static_route` - (Optional) One or more `static_route` blocks as defined below.
* `switch_id` - (Optional) . Changing this forces a new resource to be created.
* `syslog_host` - (Optional) .
* `tags` - (Optional) .
* `user` - (Optional) One or more `user` blocks as defined below.
* `vip` - (Optional) . Changing this forces a new resource to be created.
* `vrid` - (Optional) . Changing this forces a new resource to be created.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `dhcp_server` block supports the following:

* `dns_servers` - (Optional) .
* `interface_index` - (Required) .
* `range_start` - (Required) .
* `range_stop` - (Required) .

---

A `dhcp_static_mapping` block supports the following:

* `ip_address` - (Required) .
* `mac_address` - (Required) .

---

A `firewall` block supports the following:

* `direction` - (Required) .
* `expression` - (Required) One or more `expression` blocks as defined below.
* `interface_index` - (Optional) .

---

A `expression` block supports the following:

* `allow` - (Required) .
* `description` - (Optional) .
* `destination_network` - (Optional) .
* `destination_port` - (Optional) .
* `logging` - (Optional) .
* `protocol` - (Required) .
* `source_network` - (Optional) .
* `source_port` - (Optional) .

---

A `l2tp` block supports the following:

* `pre_shared_secret` - (Required) .
* `range_start` - (Required) .
* `range_stop` - (Required) .

---

A `network_interface` block supports the following:

* `index` - (Required) .
* `ip_addresses` - (Required) .
* `netmask` - (Required) .
* `switch_id` - (Required) .
* `vip` - (Optional) .

---

A `port_forwarding` block supports the following:

* `description` - (Optional) .
* `private_ip` - (Required) .
* `private_port` - (Required) .
* `protocol` - (Required) .
* `public_port` - (Required) .

---

A `pptp` block supports the following:

* `range_start` - (Required) .
* `range_stop` - (Required) .

---

A `site_to_site_vpn` block supports the following:

* `local_prefix` - (Required) .
* `peer` - (Required) .
* `pre_shared_secret` - (Required) .
* `remote_id` - (Required) .
* `routes` - (Required) .

---

A `static_nat` block supports the following:

* `description` - (Optional) .
* `private_ip` - (Required) .
* `public_ip` - (Required) .

---

A `static_route` block supports the following:

* `next_hop` - (Required) .
* `prefix` - (Required) .

---

A `user` block supports the following:

* `name` - (Required) .
* `password` - (Required) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the VPCRouter

* `read` -   (Defaults to 5 minutes) Used when reading the VPCRouter

* `update` - (Defaults to 60 minutes) Used when updating the VPCRouter

* `delete` - (Defaults to 20 minutes) Used when deregistering VPCRouter



## Attribute Reference

* `id` - The ID of the VPCRouter.
* `public_ip` - .




