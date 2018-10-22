---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_vpc_router"
sidebar_current: "docs-sakuracloud-datasource-vpc-router"
description: |-
  Get information on a SakuraCloud VPC Router.
---

# sakuracloud\_vpc\_router

Use this data source to retrieve information about a SakuraCloud VPC Router.

## Example Usage

```hcl
data "sakuracloud_vpc_router" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.
 * `zone` - (Optional) The ID of the zone.

## Attributes Reference

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

