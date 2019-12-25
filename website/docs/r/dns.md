---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns"
subcategory: "Global"
description: |-
  Manages a SakuraCloud DNS.
---

# sakuracloud_dns

Manages a SakuraCloud DNS.

## Argument Reference

* `description` - (Optional) .
* `icon_id` - (Optional) .
* `record` - (Optional) One or more `record` blocks as defined below.
* `tags` - (Optional) .
* `zone` - (Required) . Changing this forces a new resource to be created.


---

A `record` block supports the following:

* `name` - (Required) .
* `port` - (Optional) .
* `priority` - (Optional) .
* `ttl` - (Optional) .
* `type` - (Required) .
* `value` - (Required) .
* `weight` - (Optional) .


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the DNS

* `read` -   (Defaults to 5 minutes) Used when reading the DNS

* `update` - (Defaults to 5 minutes) Used when updating the DNS

* `delete` - (Defaults to 5 minutes) Used when deregistering DNS



## Attribute Reference

* `id` - The ID of the DNS.
* `dns_servers` - .




