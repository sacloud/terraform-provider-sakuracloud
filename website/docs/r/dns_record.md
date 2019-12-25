---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns_record"
subcategory: "Global"
description: |-
  Manages a SakuraCloud DNS Record.
---

# sakuracloud_dns_record

Manages a SakuraCloud DNS Record.

## Argument Reference

* `dns_id` - (Required) . Changing this forces a new resource to be created.
* `name` - (Required) . Changing this forces a new resource to be created.
* `port` - (Optional) . Changing this forces a new resource to be created.
* `priority` - (Optional) . Changing this forces a new resource to be created.
* `ttl` - (Optional) . Changing this forces a new resource to be created. Defaults to `3600`.
* `type` - (Required) . Changing this forces a new resource to be created.
* `value` - (Required) . Changing this forces a new resource to be created.
* `weight` - (Optional) . Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the DNS Record

* `read` -   (Defaults to 5 minutes) Used when reading the DNS Record


* `delete` - (Defaults to 5 minutes) Used when deregistering DNS Record



## Attribute Reference

* `id` - The ID of the DNS Record.




