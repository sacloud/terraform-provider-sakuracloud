---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_sim"
subcategory: "SecureMobile"
description: |-
  Manages a SakuraCloud SIM.
---

# sakuracloud_sim

Manages a SakuraCloud SIM.

## Argument Reference

* `carrier` - (Required) .
* `description` - (Optional) .
* `enabled` - (Optional) . Defaults to `true`.
* `iccid` - (Required) . Changing this forces a new resource to be created.
* `icon_id` - (Optional) .
* `imei` - (Optional) .
* `name` - (Required) .
* `passcode` - (Required) . Changing this forces a new resource to be created.
* `tags` - (Optional) .



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the SIM

* `read` -   (Defaults to 5 minutes) Used when reading the SIM

* `update` - (Defaults to 5 minutes) Used when updating the SIM

* `delete` - (Defaults to 5 minutes) Used when deregistering SIM



## Attribute Reference

* `id` - The id of the SIM.
* `ip_address` - .
* `mobile_gateway_id` - .




