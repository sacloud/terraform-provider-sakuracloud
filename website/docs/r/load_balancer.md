---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
subcategory: "Appliance"
description: |-
  Manages a SakuraCloud Load Balancer.
---

# sakuracloud_load_balancer

Manages a SakuraCloud Load Balancer.

## Example Usage

```hcl
resource "sakuracloud_load_balancer" "foobar" {
  name         = "foobar"
  switch_id    = sakuracloud_switch.foobar.id
  vrid         = 1
  ip_addresses = ["192.168.11.101"]
  netmask      = 24
  gateway      = "192.168.11.1"

  description = "description"
  tags        = ["tag1", "tag2"]

  vip {
    vip          = "192.168.11.201"
    port         = 80
    delay_loop   = 10
    sorry_server = "192.168.11.21"

    server {
      ip_address     = "192.168.11.51"
      check_protocol = "http"
      check_path     = "/health"
      check_status   = 200
    }

    server {
      ip_address     = "192.168.11.52"
      check_protocol = "http"
      check_path     = "/health"
      check_status   = 200
    }
  }
}

resource "sakuracloud_switch" "foobar" {
  name = "foobar"
}
```
## Argument Reference

* `description` - (Optional) The description of the LoadBalancer. The length of this value must be in the range [`1`-`512`].
* `gateway` - (Optional) The IP address of the gateway used by LoadBalancer. Changing this forces a new resource to be created.
* `high_availability` - (Optional) The flag to enable HA mode. Changing this forces a new resource to be created.
* `icon_id` - (Optional) The icon id to attach to the LoadBalancer.
* `ip_addresses` - (Required) A list of IP address to assign to the LoadBalancer. . Changing this forces a new resource to be created.
* `name` - (Required) The name of the LoadBalancer. The length of this value must be in the range [`1`-`64`].
* `netmask` - (Required) The bit length of the subnet assigned to the LoadBalancer. This must be in the range [`8`-`29`]. Changing this forces a new resource to be created.
* `plan` - (Optional) The plan name of the LoadBalancer. This must be one of [`standard`/`highspec`]. Changing this forces a new resource to be created. Default:`standard`.
* `switch_id` - (Required) The id of the switch to which the LoadBalancer connects. Changing this forces a new resource to be created.
* `tags` - (Optional) Any tags to assign to the LoadBalancer.
* `vip` - (Optional) One or more `vip` blocks as defined below.
* `vrid` - (Required) The Virtual Router Identifier. This is only used when `high_availability` is set `true`. Changing this forces a new resource to be created.
* `zone` - (Optional) The name of zone that the LoadBalancer will be created. (e.g. `is1a`, `tk1a`). Changing this forces a new resource to be created.


---

A `vip` block supports the following:

* `delay_loop` - (Optional) The interval in seconds between checks. This must be in the range [`10`-`2147483647`].
* `description` - (Optional) The description of the VIP. The length of this value must be in the range [`1`-`512`].
* `port` - (Required) The target port number for load-balancing. This must be in the range [`1`-`65535`].
* `server` - (Optional) One or more `server` blocks as defined below.
* `sorry_server` - (Optional) The IP address of the SorryServer. This will be used when all servers under this VIP are down.
* `vip` - (Required) The virtual IP address.

---

A `server` block supports the following:

* `check_path` - (Optional) The path used when checking by HTTP/HTTPS.
* `check_protocol` - (Required) The protocol used for health checks. This must be one of [`http`/`https`/`tcp`/`ping`].
* `check_status` - (Optional) The response code to expect when checking by HTTP/HTTPS.
* `enabled` - (Optional) The flag to enable as destination of load balancing.
* `ip_address` - (Required) The IP address of the destination server.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Load Balancer


* `update` - (Defaults to 60 minutes) Used when updating the Load Balancer

* `delete` - (Defaults to 20 minutes) Used when deregistering Load Balancer



## Attribute Reference

* `id` - The id of the Load Balancer.




