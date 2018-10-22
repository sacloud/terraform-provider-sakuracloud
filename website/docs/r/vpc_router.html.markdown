---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router"
sidebar_current: "docs-sakuracloud-resource-vpc-router"
description: |-
  Provides a SakuraCloud VPC Router resource. This can be used to create, update, and delete VPC Routers.
---

# sakuracloud\_vpc\_router

Provides a SakuraCloud VPC Router resource. This can be used to create, update, and delete VPC Routers.

## Example Usage

```hcl
# Create a new VPC Router(standard)
resource "sakuracloud_vpc_router" "foobar" {
  name = "foobar"

  #syslog_host         = "192.168.11.1"
  #internet_connection = true

  description = "description"
  tags        = ["foo", "bar"]
}

# Create a new VPC Router(premium or highspec)
resource "sakuracloud_vpc_router" "foobar1" {
  name       = "foobar"
  plan       = "premium"
  switch_id  = sakuracloud_internet.foobar.switch_id
  vip        = sakuracloud_internet.foobar.ipaddresses[0]
  ipaddress1 = sakuracloud_internet.foobar.ipaddresses[1]
  ipaddress2 = sakuracloud_internet.foobar.ipaddresses[2]
  #aliases   = [sakuracloud_internet.foobar.ipaddresses[3]] 
  vrid = 1
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `plan` - (Optional) The plan of the VPC Router.   
Valid value is one of the following: [ "standard" (default) / "premium" / "highspec" ]
* `switch_id` - (Required) The ID of the switch connected to the VPC Router. Used when plan is `premium` or `highspec`.
* `vrid` - (Required) VRID used when plan is `premium` or `highspec`.
* `ipaddress1` - (Required) The primary IP address of the VPC Router.
* `ipaddress2` - (Optional) The secondly IP address of the VPC Router. Used when plan is `premium` or `highspec`.
* `vip` - (Optional) The Virtual IP address of the VPC Router. Used when plan is `premium` or `highspec`.
* `aliases` - (Optional) The IP address aliase list. Used when plan is `premium` or `highspec`.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.
* `graceful_shutdown_timeout` - (Optional) The wait time (seconds) to do graceful shutdown the VPC Router.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - The name of the resource.
* `plan` - The name of the resource plan. 
* `switch_id` - The ID of the Switch connected to the VPC Router (eth0).
* `vip` - Virtual IP address of the VPC Router. Used when plan is in `premium` or `highspec`.
* `ipaddress1` - The primary IP address of the VPC Router.
* `ipaddress2` - The secondly IP address of the VPC Router. Used when plan is in `premium` or `highspec`.
* `vrid` - VRID used when plan is in `premium` or `highspec`.
* `aliases` - The IP address aliase list. Used when plan is in `premium` or `highspec`.
* `global_address` - Global IP address of the VPC Router.
* `syslog_host` - The destination HostName/IP address to send log.	
* `internet_connection` - The flag of enable/disable connection from the VPC Router to the Internet.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import

VPC Routers can be imported using the VPC Router ID.

```
$ terraform import sakuracloud_vpc_router.foobar <vpc_router_id>
```
