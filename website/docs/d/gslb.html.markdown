---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_gslb"
sidebar_current: "docs-sakuracloud-datasource-gslb"
description: |-
  Get information on a SakuraCloud GSLB.
---

# sakuracloud\_gslb

Use this data source to retrieve information about a SakuraCloud GSLB.

## Example Usage

```hcl
data "sakuracloud_gslb" "foobar" {
  name_selectors = ["foobar"]
}
```

## Argument Reference

 * `name_selectors` - (Optional) The list of names to filtering.
 * `tag_selectors` - (Optional) The list of tags to filtering.
 * `filter` - (Optional) The map of filter key and value.

## Attributes Reference

* `id` - The ID of the resource.
* `name` - Name of the resource.
* `fqdn` - FQDN to access this resource.
* `health_check` - Health check rules. It contains some attributes to [Health Check](#health-check).
* `weighted` - The flag for enable/disable weighting.
* `sorry_server` - The hostname or IP address of sorry server.
* `description` - The description of the resource.
* `tags` - The tag list of the resources.
* `icon_id` - The ID of the icon of the resource.

### Health Check

Attributes for Health Check:

* `protocol` - Protocol used in health check.
* `delay_loop` - Health check access interval (unit:`second`).
* `host_header` - The value of `Host` header used in http/https health check access.
* `path` - The request path used in http/https health check access.
* `status` - HTTP status code expected by health check access.
* `port` - Port number used in tcp health check access.
