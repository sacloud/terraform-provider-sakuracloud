---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_gslb"
sidebar_current: "docs-sakuracloud-resource-global-gslb-setting"
description: |-
  Provides a SakuraCloud GSLB resource. This can be used to create, update, and delete GSLBs.
---

# sakuracloud\_gslb

Provides a SakuraCloud GSLB resource. This can be used to create, update, and delete GSLBs.

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

  sorry_server = "192.2.0.1"

  description = "description"
  tags        = ["foo", "bar"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the resource.
* `health_check` - (Required) Health check rules. It contains some attributes to [Health Check](#health-check).
* `weighted` - (Optional) The flag for enable/disable weighting (default:`true`).
* `sorry_server` - (Optional) The hostname or IP address of sorry server.
* `description` - (Optional) The description of the resource.
* `tags` - (Optional) The tag list of the resources.
* `icon_id` - (Optional) The ID of the icon.

### Health Check

Attributes for Health Check:

* `protocol` - (Required) Protocol used in health check.  
Valid value is one of the following: [ "http" / "https" / "ping" / "tcp" ]
* `delay_loop` - (Optional) Health check access interval (unit:`second`, default:`10`).
* `host_header` - (Optional) The value of `Host` header used in http/https health check access.
* `path` - (Optional) The request path used in http/https health check access.
* `status` - (Optional) HTTP status code expected by health check access.
* `port` - (Optional) Port number used in tcp health check access.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `name` - Name of the resource.
* `fqdn` - FQDN to access this resource.
* `health_check` - Health check rules. It contains some attributes to [Health Check](#health-check).
* `weighted` - The flag for enable/disable weighting.
* `sorry_server` - The hostname or IP address of sorry server.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.

## Import

GSLBs can be imported using the GSLB ID.

```
$ terraform import sakuracloud_gslb.foobar <gslb_id>
```
