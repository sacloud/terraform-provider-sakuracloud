---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_mobile_gateway"
subcategory: "SecureMobile"
description: |-
  Manages a SakuraCloud MobileGateway.
---

# sakuracloud_mobile_gateway

Manages a SakuraCloud MobileGateway.

## Argument Reference

* `description` - (Optional) .
* `dns_servers` - (Required) .
* `icon_id` - (Optional) .
* `inter_device_communication` - (Optional) .
* `internet_connection` - (Optional) .
* `name` - (Required) .
* `private_network_interface` - (Optional) An `private_network_interface` block as defined below.
* `sim` - (Optional) One or more `sim` blocks as defined below.
* `sim_route` - (Optional) One or more `sim_route` blocks as defined below.
* `static_route` - (Optional) One or more `static_route` blocks as defined below.
* `tags` - (Optional) .
* `traffic_control` - (Optional) A `traffic_control` block as defined below.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `private_network_interface` block supports the following:

* `ip_address` - (Required) .
* `netmask` - (Required) .
* `switch_id` - (Required) .

---

A `sim` block supports the following:

* `ip_address` - (Required) .
* `sim_id` - (Required) .

---

A `sim_route` block supports the following:

* `prefix` - (Required) .
* `sim_id` - (Required) .

---

A `static_route` block supports the following:

* `next_hop` - (Required) .
* `prefix` - (Required) .

---

A `traffic_control` block supports the following:

* `auto_traffic_shaping` - (Optional) .
* `band_width_limit` - (Optional) .
* `enable_email` - (Optional) .
* `enable_slack` - (Optional) .
* `quota` - (Required) .
* `slack_webhook` - (Optional) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the MobileGateway

* `read` -   (Defaults to 5 minutes) Used when reading the MobileGateway

* `update` - (Defaults to 60 minutes) Used when updating the MobileGateway

* `delete` - (Defaults to 20 minutes) Used when deregistering MobileGateway



## Attribute Reference

* `id` - The ID of the MobileGateway.
* `public_ip` - .
* `public_netmask` - .




