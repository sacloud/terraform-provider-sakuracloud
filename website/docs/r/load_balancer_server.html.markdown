---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer_server"
sidebar_current: "docs-sakuracloud-resource-lb-server"
description: |-
  Provides a SakuraCloud LoadBalancer Server resource. This can be used to create and delete LoadBalancer Servers.
---

# sakuracloud\_load\_balancer\_server

Provides a SakuraCloud LoadBalancer Server resource. This can be used to create and delete LoadBalancer Servers.

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
}

# Create new LoadBalancer servers
resource "sakuracloud_load_balancer_server" "server01" {
  load_balancer_vip_id = sakuracloud_load_balancer_vip.vip1.id
  ipaddress            = "192.168.2.151"
  check_protocol       = "https"
  check_path           = "/healthz"
  check_status         = 200
  enabled              = true
}

resource "sakuracloud_load_balancer_server" "server02" {
  load_balancer_vip_id = sakuracloud_load_balancer_vip.vip1.id
  ipaddress            = "192.168.2.152"
  check_protocol       = "https"
  check_path           = "/healthz"
  check_status         = 200
  enabled              = true
}



```

## Argument Reference

The following arguments are supported:

* `load_vip_id` - (Required) The ID of the Load Balancer VIP to which the Server belongs.
* `ipaddress` - (Required) The IP address of the Server.
* `check_protocol` - (Required) Protocol used in health check.  
Valid value is one of the following: [ "http" / "https" / "ping" / "tcp" ]
* `check_path` - (Optional) The request path used in http/https health check access.
* `check_status` - (Optional) HTTP status code expected by health check access.
* `enabled` - (Optional) The flag of enable/disable the Server.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `ipaddress` - The IP address of the Server.
* `check_protocol` - Protocol used in health check.
* `check_path` - The request path used in http/https health check access.
* `check_status` - HTTP status code expected by health check access.
* `enabled` - The flag of enable/disable the Server.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of Load Balancer Server is not supported.

