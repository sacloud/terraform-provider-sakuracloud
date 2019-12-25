---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router"
subcategory: "Appliance"
description: |-
  Get information about an existing VPCRouter.
---

# Data Source: sakuracloud_vpc_router

Get information about an existing VPCRouter.

## Argument Reference

* `filter` - (Optional) A `filter` block as defined below.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `filter` block supports the following:

* `condition` - (Optional) One or more `condition` blocks as defined below.
* `id` - (Optional) .
* `names` - (Optional) .
* `tags` - (Optional) .

---

A `condition` block supports the following:

* `name` - (Required) .
* `values` - (Required) .


## Attribute Reference

* `id` - The ID of the VPCRouter.
* `aliases` - .
* `description` - .
* `dhcp_server` - A list of `dhcp_server` blocks as defined below.
* `dhcp_static_mapping` - A list of `dhcp_static_mapping` blocks as defined below.
* `firewall` - A list of `firewall` blocks as defined below.
* `icon_id` - .
* `internet_connection` - .
* `ip_addresses` - .
* `l2tp` - A list of `l2tp` blocks as defined below.
* `name` - .
* `network_interface` - A list of `network_interface` blocks as defined below.
* `plan` - .
* `port_forwarding` - A list of `port_forwarding` blocks as defined below.
* `pptp` - A list of `pptp` blocks as defined below.
* `public_ip` - .
* `site_to_site_vpn` - A list of `site_to_site_vpn` blocks as defined below.
* `static_nat` - A list of `static_nat` blocks as defined below.
* `static_route` - A list of `static_route` blocks as defined below.
* `switch_id` - .
* `syslog_host` - .
* `tags` - .
* `user` - A list of `user` blocks as defined below.
* `vip` - .
* `vrid` - .


---

A `dhcp_server` block exports the following:

* `dns_servers` - .
* `interface_index` - .
* `range_start` - .
* `range_stop` - .

---

A `dhcp_static_mapping` block exports the following:

* `ip_address` - .
* `mac_address` - .

---

A `firewall` block exports the following:

* `direction` - .
* `expression` - A list of `expression` blocks as defined below.
* `interface_index` - .

---

A `expression` block exports the following:

* `allow` - .
* `description` - .
* `dest_nw` - .
* `dest_port` - .
* `logging` - .
* `protocol` - .
* `source_nw` - .
* `source_port` - .

---

A `l2tp` block exports the following:

* `pre_shared_secret` - .
* `range_start` - .
* `range_stop` - .

---

A `network_interface` block exports the following:

* `index` - .
* `ip_addresses` - .
* `netmask` - .
* `switch_id` - .
* `vip` - .

---

A `port_forwarding` block exports the following:

* `description` - .
* `private_ip` - .
* `private_port` - .
* `protocol` - .
* `public_port` - .

---

A `pptp` block exports the following:

* `range_start` - .
* `range_stop` - .

---

A `site_to_site_vpn` block exports the following:

* `local_prefix` - .
* `peer` - .
* `pre_shared_secret` - .
* `remote_id` - .
* `routes` - .

---

A `static_nat` block exports the following:

* `description` - .
* `private_ip` - .
* `public_ip` - .

---

A `static_route` block exports the following:

* `next_hop` - .
* `prefix` - .

---

A `user` block exports the following:

* `name` - .
* `password` - .



