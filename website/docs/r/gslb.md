---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_gslb"
subcategory: "Global"
description: |-
  Manages a SakuraCloud GSLB.
---

# sakuracloud_gslb

Manages a SakuraCloud GSLB.

## Example Usage

```hcl
resource "sakuracloud_gslb" "foobar" {
  name = "example"

  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "example.com"
    path        = "/"
    status      = "200"
  }

  sorry_server = "192.2.0.1"

  server {
    ip_address = "192.2.0.11"
    weight     = 1
    enabled    = true
  }
  server {
    ip_address = "192.2.0.12"
    weight     = 1
    enabled    = true
  }

  description = "description"
  tags        = ["tag1", "tag2"]
}
```
## Argument Reference

* `description` - (Optional) The description of the GSLB. The length of this value must be in the range [`1`-`512`].
* `health_check` - (Required) A `health_check` block as defined below.
* `icon_id` - (Optional) The icon id to attach to the GSLB.
* `name` - (Required) The name of the GSLB. The length of this value must be in the range [`1`-`64`].
* `server` - (Optional) One or more `server` blocks as defined below.
* `sorry_server` - (Optional) The IP address of the SorryServer. This will be used when all servers are down.
* `tags` - (Optional) Any tags to assign to the GSLB.
* `weighted` - (Optional) The flag to enable weighted load-balancing.


---

A `health_check` block supports the following:

* `delay_loop` - (Optional) The interval in seconds between checks. This must be in the range [`10`-`60`].
* `host_header` - (Optional) The value of host header send when checking by HTTP/HTTPS.
* `path` - (Optional) The path used when checking by HTTP/HTTPS.
* `port` - (Optional) The port number used when checking by TCP.
* `protocol` - (Required) The protocol used for health checks. This must be one of [`http`/`https`/`tcp`/`ping`].
* `status` - (Optional) The response-code to expect when checking by HTTP/HTTPS.

---

A `server` block supports the following:

* `enabled` - (Optional) The flag to enable as destination of load balancing.
* `ip_address` - (Required) The IP address of the server.
* `weight` - (Optional) The weight used when weighted load balancing is enabled. This must be in the range [`1`-`10000`].


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the GSLB


* `update` - (Defaults to 5 minutes) Used when updating the GSLB

* `delete` - (Defaults to 5 minutes) Used when deregistering GSLB



## Attribute Reference

* `id` - The id of the GSLB.
* `fqdn` - The FQDN for accessing to the GSLB. This is typically used as value of CNAME record.




