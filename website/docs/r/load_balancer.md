---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
subcategory: "Appliance"
description: |-
  Provides a SakuraCloud LoadBalancer resource. This can be used to create, update, and delete LoadBalancers.
---

# sakuracloud\_load\_balancer

Provides a SakuraCloud LoadBalancer resource. This can be used to create, update, and delete LoadBalancers.

## Example Usage

```hcl
# Create a new LoadBalancer
resource "sakuracloud_load_balancer" "foobar" {
  name = "foobar"

  switch_id         = sakuracloud_switch.foobar.id
  vrid              = 1
  high_availability = false
  plan              = "standard"

  ipaddress1 = "192.168.2.1"
  # only when high_availability = true 
  #ipaddress2        = "192.168.2.2"
  
  nw_mask_len   = 24
  default_route = "192.168.2.254"

  vips {
    vip          = "192.168.2.101"
    port         = 80
    delay_loop   = 10
    sorry_server = "192.168.2.11"
    
    servers {
      ipaddress      = "192.168.2.102"
      check_protocol = "http"
      check_path     = "/ping.html"
      check_status   = 200
    }
    servers {
      ipaddress      = "192.168.2.103"
      check_protocol = "http"
      check_path     = "/ping.html"
      check_status   = 200
    }
  }

  description = "description"
  tags        = ["foo", "bar"]
}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `switch_id` - (Required) The ID of the Switch connected to the Load Balancer.
* `vrid` - (Required) VRID used when high-availability mode enabled.
* `high_availability` - (Optional) The flag of enable/disable high-availability mode.
* `plan` - (Optional) The name of the resource plan.  
Valid value is one of the following: [ "standard" (default) / "highspec"]
* `ipaddress1` - (Required) The primary IP address of the Load Balancer.
* `ipaddress2` - (Optional) The secondly IP address of the Load Balancer. Used when high-availability mode enabled.
* `nw_mask_len` - (Required) Network mask length.
* `default_route` - (Optional) Default gateway address of the Load Balancer.	 
* `vips` - (Optional) VIPs. It contains some attributes to [VIPs](#vips).
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

### VIPs

Attributes for VIPs:

* `vip` - (Required) The virtual IP address.
* `port` - (Required) The port number on which Load Balancer listens.
* `delay_loop` - (Optional) The interval seconds for health check access.
* `sorry_server` - (Optional) The hostname or IP address of sorry server.
* `description` - (Optional) The description of the VIP.
* `servers` - (Optional) Real servers. It contains some attributes to [Servers](#servers).

### Servers

Attributes for Servers:

* `ipaddress` - (Required) The IP address of the Server.
* `check_protocol` - (Required) Protocol used in health check.  
Valid value is one of the following: [ "http" / "https" / "ping" / "tcp" ]
* `check_path` - (Optional) The request path used in http/https health check access.
* `check_status` - (Optional) HTTP status code expected by health check access.
* `enabled` - (Optional) The flag of enable/disable the Server.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `switch_id` - The ID of the Switch connected to the Load Balancer.
* `vrid` - VRID used when high-availability mode enabled.
* `high_availability` - The flag of enable/disable high-availability mode.
* `plan` - The name of the resource plan. 
* `ipaddress1` - The primary IP address of the Load Balancer.
* `ipaddress2` - The secondly IP address of the Load Balancer. Used when high-availability mode enabled.
* `nw_mask_len` - Network mask length.
* `default_route` - Default gateway address of the Load Balancer.	 
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

LoadBalancers can be imported using the LoadBalancer ID.

```
$ terraform import sakuracloud_load_balancer.foobar <load_balancer_id>
```
