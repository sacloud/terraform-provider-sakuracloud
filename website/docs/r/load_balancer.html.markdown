---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
sidebar_current: "docs-sakuracloud-resource-lb-load-balancer"
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
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the server connected to the resource.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

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
