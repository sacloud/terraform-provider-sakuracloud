---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_gslb_server"
sidebar_current: "docs-sakuracloud-resource-global-gslb-server"
description: |-
  Provides a SakuraCloud GSLB Server resource. This can be used to create and delete GSLB Servers.
---

# sakuracloud\_gslb

Provides a SakuraCloud GSLB Server resource. This can be used to create and delete GSLB Servers.

## Example Usage

```hcl
# Create a new GSLB
resource "sakuracloud_gslb" "foobar" {
  name = "foobar"
  health_check {
    protocol    = "https"
    delay_loop  = 20
    host_header = "example.com"
    path        = "/"
    status      = "200"
  }
}

# Add Server To GSLB
resource "sakuracloud_gslb_server" "foobar1" {
  gslb_id   = sakuracloud_gslb.foobar.id
  ipaddress = "192.2.0.1"
  enabled   = true
  weight    = 1
}

# Add Server To GSLB
resource "sakuracloud_gslb_server" "foobar2" {
  gslb_id   = sakuracloud_gslb.foobar.id
  ipaddress = "192.2.0.2"
  enabled   = true
  weight    = 1
}

```

## Argument Reference

The following arguments are supported:

* `gslb_id` - (Required) The ID of the GSLB to which the GSLB Server belongs.
* `ipaddress` - (Required) The IP address of the GSLB Server.
* `enabled` - (Optional) The flag for enable/disable the GSLB Server (default:`true`).
* `weight` - (Optional) The weight of GSLB server used when weighting is enabled in the GSLB.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `ipaddress` - The IP address of the GSLB Server.
* `enabled` - The flag for enable/disable the GSLB Server.
* `weight` - The weight of GSLB server used when weighting is enabled in the GSLB.

## Import (not supported)

Import of GSLB Server is not supported.
