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

* `dns_id` - (Required) The id of the DNS resource. Changing this forces a new resource to be created.
* `name` - (Required) The name of the DNS Record resource. Changing this forces a new resource to be created.
* `port` - (Optional) The number of port. This must be in the range [`1`-`65535`]. Changing this forces a new resource to be created.
* `priority` - (Optional) The priority of target DNS Record. This must be in the range [`0`-`65535`]. Changing this forces a new resource to be created.
* `ttl` - (Optional) The number of the TTL. Changing this forces a new resource to be created. Default:`3600`.
* `type` - (Required) The type of DNS Record. This must be one of [`A`/`AAAA`/`ALIAS`/`CNAME`/`NS`/`MX`/`TXT`/`SRV`/`CAA`/`PTR`]. Changing this forces a new resource to be created.
* `value` - (Required) The value of the DNS Record. Changing this forces a new resource to be created.
* `weight` - (Optional) The weight of target DNS Record. This must be in the range [`0`-`65535`]. Changing this forces a new resource to be created.



### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the DNS Record

* `read` -   (Defaults to 5 minutes) Used when reading the DNS Record


* `delete` - (Defaults to 5 minutes) Used when deregistering DNS Record



## Attribute Reference

* `id` - The id of the DNS Record.




