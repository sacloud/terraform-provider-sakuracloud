---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer_vip"
sidebar_current: "docs-sakuracloud-resource-lb-vip"
description: |-
  Provides a SakuraCloud LoadBalancer VIP resource. This can be used to create, update, and delete LoadBalancer VIPs.
---

# sakuracloud\_load\_balancer\_vip

Provides a SakuraCloud LoadBalancer VIP resource. This can be used to create, update, and delete LoadBalancer VIPs.

## Example Usage

```hcl
# Create a new LoadBalancer
resource "sakuracloud_load_balancer" "foobar" {
  name        = "foobar"
  switch_id   = sakuracloud_switch.sw.id
  vrid        = 1
  ipaddress1  = "192.168.2.1"
  nw_mask_len = 24
}

# Create a new LoadBalancer VIP
resource "sakuracloud_load_balancer_vip" "vip1" {
  load_balancer_id = sakuracloud_load_balancer.foobar.id
  vip              = "192.168.2.101"
  port             = 80
  delay_loop       = 50
  sorry_server     = "192.168.2.201"
  description      = "description"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of the Load Balancer to which the VIP belongs.
* `vip` - (Required) The virtual IP address.
* `port` - (Required) The port number on which Load Balancer listens.
* `delay_loop` - (Optional) The interval seconds for health check access.
* `sorry_server` - (Optional) The hostname or IP address of sorry server.
* `description` - (Optional) The description of the VIP.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `vip` - The virtual IP address.
* `port` - The port number on which Load Balancer listens.
* `delay_loop` - The interval seconds for Health check access.
* `sorry_server` - The hostname or IP address of sorry server.
* `description` - The description of the VIP.
* `servers` - The internal ID list of servers under VIP.
* `zone` - The ID of the zone to which the resource belongs.


## Import (not supported)

Import of Load Balancer VIP is not supported.

