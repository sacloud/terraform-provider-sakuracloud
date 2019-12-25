---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_load_balancer"
description: |-
  Manages a SakuraCloud sakuracloud_load_balancer.
---

# sakuracloud_load_balancer

Manages a SakuraCloud sakuracloud_load_balancer.

## Argument Reference

* `description` - (Optional) .
* `gateway` - (Optional) . Changing this forces a new resource to be created.
* `high_availability` - (Optional) . Changing this forces a new resource to be created.
* `icon_id` - (Optional) .
* `ip_addresses` - (Required) . Changing this forces a new resource to be created.
* `name` - (Required) .
* `netmask` - (Required) . Changing this forces a new resource to be created.
* `plan` - (Optional) . Changing this forces a new resource to be created. Defaults to `standard`.
* `switch_id` - (Required) . Changing this forces a new resource to be created.
* `tags` - (Optional) .
* `vip` - (Optional) One or more `vip` blocks as defined below.
* `vrid` - (Required) . Changing this forces a new resource to be created.
* `zone` - (Optional) target SakuraCloud zone. Changing this forces a new resource to be created.


---

A `vip` block supports the following:

* `delay_loop` - (Optional) .
* `description` - (Optional) .
* `port` - (Required) .
* `server` - (Optional) One or more `server` blocks as defined below.
* `sorry_server` - (Optional) .
* `vip` - (Required) .

---

A `server` block supports the following:

* `check_path` - (Optional) .
* `check_protocol` - (Required) .
* `check_status` - (Optional) .
* `enabled` - (Optional) .
* `ip_address` - (Required) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the sakuracloud_load_balancer

* `read` -   (Defaults to 5 minutes) Used when reading the sakuracloud_load_balancer

* `update` - (Defaults to 60 minutes) Used when updating the sakuracloud_load_balancer

* `delete` - (Defaults to 20 minutes) Used when deregistering sakuracloud_load_balancer



## Attribute Reference

* `id` - The ID of the sakuracloud_load_balancer.




