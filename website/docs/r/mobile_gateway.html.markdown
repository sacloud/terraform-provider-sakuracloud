---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_mobile_gateway"
sidebar_current: "docs-sakuracloud-resource-secure-mobile-mobile-gateway"
description: |-
  Provides a SakuraCloud Mobile Gateway resource. This can be used to create, update, and delete Mobile Gateways.
---

# sakuracloud\_mobile\_gateway

Provides a SakuraCloud Mobile Gateway resource. This can be used to create, update, and delete Mobile Gateways.

## Example Usage

```hcl
# Create a new Mobile Gateway
resource "sakuracloud_mobile_gateway" "foobar" {
  name = "foobar"

  switch_id           = sakuracloud_switch.sw.id
  private_ipaddress   = "192.168.11.101"
  private_nw_mask_len = 24
  internet_connection = true
  dns_server1         = "8.8.8.8"
  dns_server2         = "8.8.4.4"

  traffic_control = {
    quota                = 256
    band_width_limit     = 64
    enable_email         = true
    enable_slack         = true
    slack_webhook        = "https://hooks.slack.com/services/xxx/xxx/xxx"
    auto_traffic_shaping = true
  }

  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `switch_id` - (Optional) The ID of the switch connected to the Mobile Gateway.
* `private_ipaddress` - (Optional) The IP address on private NIC of the Mobile Gateway.
* `private_nw_mask_len` - (Optional) The network mask length on private NIC of the Mobile Gateway.
* `internet_connection` - (Optional) The flag of enable/disable connecting from MobileGateway to the Internet.
* `dns_server1` - (Optional) The primary DNS server IP address.
* `dns_server2` - (Optional) The secondly DNS server IP address.
* `traffic_control` - (Optional) Traffic control rules. It contains some attributes to [Traffic Control](#traffic-control).
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

### Traffic Control

Attributes for Traffic Control:

* `quota` - (Required) Traffic quota size (unit:`MB`).  
* `band_width_limit` - (Optional) Traffic bandwidth limit(unit:`Kbps`). 
* `auto_traffic_shaping` - (Optional) The flag of enable/disable Auto Traffic Shaping.
* `enable_email` - (Optional) The flag of enable/disable e-mail notification.
* `enable_slack` - (Optional) The flag of enable/disable slack notification.
* `slack_webhook` - (Optional) The webhook URL of destination of slack notification.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `switch_id` - The ID of the switch connected to the Mobile Gateway.
* `public_ipaddress` - The IP address on public NIC of the Mobile Gateway.
* `public_nw_mask_len` - The network mask length on public NIC of the Mobile Gateway.
* `private_ipaddress` - The IP address on private NIC of the Mobile Gateway.
* `private_nw_mask_len` - The network mask length on private NIC of the Mobile Gateway.
* `internet_connection` - The flag of enable/disable connecting from MobileGateway to the Internet.
* `dns_server1` - The primary DNS server IP address.
* `dns_server2` - The secondly DNS server IP address.
* `traffic_control` - Traffic control rules.
* `sim_ids` - The ID list of the SIMs connected to the Mobile Gateway.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon.
* `zone` - The ID of the zone to which the resource belongs.

## Import

Mobile Gateways can be imported using the Mobile Gateway ID.

```
$ terraform import sakuracloud_mobile_gateway.foobar <mobile_gateway_id>
```
