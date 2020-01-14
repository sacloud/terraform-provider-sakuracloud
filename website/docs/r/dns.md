---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_dns"
subcategory: "Global"
description: |-
  Manages a SakuraCloud DNS.
---

# sakuracloud_dns

Manages a SakuraCloud DNS.

## Example Usage

```hcl
resource "sakuracloud_dns" "foobar" {
  zone        = "example.com"
  description = "description"
  tags        = ["tag1", "tag2"]
  record {
    name  = "www"
    type  = "A"
    value = "192.168.11.1"
  }
  record {
    name  = "www"
    type  = "A"
    value = "192.168.11.2"
  }
}
```
## Argument Reference

* `description` - (Optional) The description of the DNS. The length of this value must be in the range [`1`-`512`].
* `icon_id` - (Optional) The icon id to attach to the DNS.
* `record` - (Optional) One or more `record` blocks as defined below.
* `tags` - (Optional) Any tags to assign to the DNS.
* `zone` - (Required) The target zone. (e.g. `example.com`). Changing this forces a new resource to be created.


---

A `record` block supports the following:

* `name` - (Required) The name of the DNS Record. The length of this value must be in the range [`1`-`64`].
* `port` - (Optional) The number of port. This must be in the range [`1`-`65535`].
* `priority` - (Optional) The priority of target DNS Record. This must be in the range [`0`-`65535`].
* `ttl` - (Optional) The number of the TTL.
* `type` - (Required) The type of DNS Record. This must be one of [`A`/`AAAA`/`ALIAS`/`CNAME`/`NS`/`MX`/`TXT`/`SRV`/`CAA`/`PTR`].
* `value` - (Required) The value of the DNS Record.
* `weight` - (Optional) The weight of target DNS Record. This must be in the range [`0`-`65535`].


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the DNS


* `update` - (Defaults to 5 minutes) Used when updating the DNS

* `delete` - (Defaults to 5 minutes) Used when deregistering DNS



## Attribute Reference

* `id` - The id of the DNS.
* `dns_servers` - A list of IP address of DNS server that manage this zone.




