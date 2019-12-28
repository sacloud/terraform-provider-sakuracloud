---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_gslb"
subcategory: "Global"
description: |-
  Manages a SakuraCloud GSLB.
---

# sakuracloud_gslb

Manages a SakuraCloud GSLB.

## Argument Reference

* `description` - (Optional) .
* `health_check` - (Required) A `health_check` block as defined below.
* `icon_id` - (Optional) .
* `name` - (Required) . Changing this forces a new resource to be created.
* `server` - (Optional) One or more `server` blocks as defined below.
* `sorry_server` - (Optional) .
* `tags` - (Optional) .
* `weighted` - (Optional) .


---

A `health_check` block supports the following:

* `delay_loop` - (Optional) .
* `host_header` - (Optional) .
* `path` - (Optional) .
* `port` - (Optional) .
* `protocol` - (Required) .
* `status` - (Optional) .

---

A `server` block supports the following:

* `enabled` - (Optional) .
* `ip_address` - (Required) .
* `weight` - (Optional) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the GSLB

* `read` -   (Defaults to 5 minutes) Used when reading the GSLB

* `update` - (Defaults to 5 minutes) Used when updating the GSLB

* `delete` - (Defaults to 5 minutes) Used when deregistering GSLB



## Attribute Reference

* `id` - The id of the GSLB.
* `fqdn` - .




