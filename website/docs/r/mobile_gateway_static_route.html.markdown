---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_mobile_gateway"
sidebar_current: "docs-sakuracloud-resource-secure-mobile-mgwstaticroute"
description: |-
  Provides a SakuraCloud Mobile Gateway Static Route resource. This can be used to create and delete Mobile Gateway Static Routes.
---

# sakuracloud\_mobile\_gateway\_static\_route

Provides a SakuraCloud Mobile Gateway Static Route resource. This can be used to create, update, and delete Mobile Gateway Static Routes.

## Example Usage

```hcl
# Create a new Mobile Gateway
resource "sakuracloud_mobile_gateway" "foobar" {
  name = "foobar"

  switch_id           = sakuracloud_switch.sw.id
  private_ipaddress   = "192.168.2.101"
  private_nw_mask_len = 24
  internet_connection = true
  dns_server1         = "8.8.8.8"
  dns_server2         = "8.8.4.4"

  description = "description"
  tags        = ["foo", "bar"]
}

# Create a new Static Route
resource "sakuracloud_mobile_gateway_static_route" "r1" {
  mobile_gateway_id = sakuracloud_mobile_gateway.foobar.id
  prefix            = "10.16.0.0/16"
  next_hop          = "192.168.2.1"
}
```

## Argument Reference

The following arguments are supported:

* `mobile_gateway_id` - (Required) The ID of the Mobile Gateway to set Static Route.
* `prefix` - (Required) The routing prefix (format:`CIDR`).
* `next_hop` - (Required) The IP address of the next hop.
* `zone` - (Optional) The ID of the zone to which the resource belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `zone` - The ID of the zone to which the resource belongs.

## Import (not supported)

Import of Mobile Gateway Static Route is not supported.
