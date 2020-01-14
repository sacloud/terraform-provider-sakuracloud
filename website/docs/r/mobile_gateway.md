---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_mobile_gateway"
subcategory: "SecureMobile"
description: |-
  Manages a SakuraCloud Mobile Gateway.
---

# sakuracloud_mobile_gateway

Manages a SakuraCloud Mobile Gateway.

## Argument Reference

* `description` - (Optional) The description of the MobileGateway. The length of this value must be in the range [`1`-`512`].
* `dns_servers` - (Required) A list of IP address used by each connected devices.
* `icon_id` - (Optional) The icon id to attach to the MobileGateway.
* `inter_device_communication` - (Optional) The flag to allow communication between each connected devices.
* `internet_connection` - (Optional) The flag to enable connect to the Internet.
* `name` - (Required) The name of the MobileGateway. The length of this value must be in the range [`1`-`64`].
* `private_network_interface` - (Optional) An `private_network_interface` block as defined below.
* `sim` - (Optional) One or more `sim` blocks as defined below.
* `sim_route` - (Optional) One or more `sim_route` blocks as defined below.
* `static_route` - (Optional) One or more `static_route` blocks as defined below.
* `tags` - (Optional) Any tags to assign to the MobileGateway.
* `traffic_control` - (Optional) A `traffic_control` block as defined below.
* `zone` - (Optional) The name of zone that the MobileGateway will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


---

A `private_network_interface` block supports the following:

* `ip_address` - (Required) The IP address to assign to the MobileGateway.
* `netmask` - (Required) The bit length of the subnet to assign to the MobileGateway. This must be in the range [`8`-`29`].
* `switch_id` - (Required) The id of the switch to which the MobileGateway connects.

---

A `sim` block supports the following:

* `ip_address` - (Required) The IP address to assign to the SIM.
* `sim_id` - (Required) The id of the Switch connected to the MobileGateway.

---

A `sim_route` block supports the following:

* `prefix` - (Required) The destination network prefix used by the sim routing. This must be specified by CIDR block formatted string.
* `sim_id` - (Required) The id of the routing destination SIM.

---

A `static_route` block supports the following:

* `next_hop` - (Required) The IP address of next hop.
* `prefix` - (Required) The destination network prefix used by static routing. This must be specified by CIDR block formatted string.

---

A `traffic_control` block supports the following:

* `auto_traffic_shaping` - (Optional) The flag to enable the traffic shaping.
* `band_width_limit` - (Optional) The bandwidth allowed when the traffic shaping is enabled.
* `enable_email` - (Optional) The flag to enable email notification when the traffic shaping is enabled.
* `enable_slack` - (Optional) The flag to enable slack notification when the traffic shaping is enabled.
* `quota` - (Required) The threshold of monthly traffic usage to enable to the traffic shaping.
* `slack_webhook` - (Optional) The webhook URL used when sends notification. It will only used when `enable_slack` is set `true`.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Mobile Gateway


* `update` - (Defaults to 60 minutes) Used when updating the Mobile Gateway

* `delete` - (Defaults to 20 minutes) Used when deregistering Mobile Gateway



## Attribute Reference

* `id` - The id of the Mobile Gateway.
* `public_ip` - The public IP address assigned to the MobileGateway.
* `public_netmask` - The bit length of the subnet assigned to the MobileGateway.




